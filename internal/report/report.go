package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/rioaguspermana/quota-race/internal/burst"
	"github.com/rioaguspermana/quota-race/internal/config"
)

type Summary struct {
	Pass     bool            `json:"pass"`
	MaxOK    int             `json:"max_ok"`
	OKStatus int             `json:"ok_status"`
	Attempts []AttemptView   `json:"attempts"`
}

type AttemptView struct {
	Index      int         `json:"index"`
	OKCount    int         `json:"ok_count"`
	ErrorCount int         `json:"error_count"`
	ByStatus   map[int]int `json:"by_status"`
	Held       bool        `json:"invariant_held"`
	ElapsedMS  int64       `json:"elapsed_ms"`
	FollowStat int         `json:"follow_up_status,omitempty"`
	FollowBody string      `json:"follow_up_body,omitempty"`
	Samples    []string    `json:"failure_samples,omitempty"`
}

func Summarize(cfg config.File, attempts []burst.Attempt) Summary {
	s := Summary{
		Pass:     true,
		MaxOK:    cfg.Invariant.MaxOK,
		OKStatus: cfg.Invariant.OKStatus,
	}
	for _, a := range attempts {
		held := a.OKCount <= cfg.Invariant.MaxOK
		if !held {
			s.Pass = false
		}
		view := AttemptView{
			Index:      a.Index,
			OKCount:    a.OKCount,
			ErrorCount: a.ErrorCount,
			ByStatus:   a.ByStatus,
			Held:       held,
			ElapsedMS:  a.Elapsed.Milliseconds(),
			FollowStat: a.FollowStat,
			FollowBody: truncate(a.FollowBody, 200),
		}
		if !held {
			view.Samples = samples(a, cfg.Invariant.OKStatus, 3)
		}
		s.Attempts = append(s.Attempts, view)
	}
	return s
}

func WriteText(w io.Writer, s Summary) {
	if s.Pass {
		fmt.Fprintln(w, "PASS  invariant held (ok responses <= max_ok on every attempt)")
	} else {
		fmt.Fprintln(w, "FAIL  invariant broken: too many successful grants under concurrency")
	}
	fmt.Fprintf(w, "rule: status %d counted as grant; max_ok=%d\n", s.OKStatus, s.MaxOK)
	for _, a := range s.Attempts {
		fmt.Fprintf(w, "  attempt %d: grants=%d errors=%d held=%v elapsed=%dms statuses=%s\n",
			a.Index, a.OKCount, a.ErrorCount, a.Held, a.ElapsedMS, formatStatuses(a.ByStatus))
		if a.FollowBody != "" {
			fmt.Fprintf(w, "    follow-up %d %s\n", a.FollowStat, a.FollowBody)
		}
		for _, sm := range a.Samples {
			fmt.Fprintf(w, "    sample: %s\n", sm)
		}
	}
}

func WriteJSON(w io.Writer, s Summary) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

func samples(a burst.Attempt, okStatus, n int) []string {
	var out []string
	for _, r := range a.Results {
		if r.Status != okStatus {
			continue
		}
		out = append(out, truncate(r.Body, 120))
		if len(out) >= n {
			break
		}
	}
	return out
}

func formatStatuses(m map[int]int) string {
	if len(m) == 0 {
		return "{}"
	}
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%d:%d", k, m[k]))
	}
	return "{" + strings.Join(parts, " ") + "}"
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
