package burst

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rioaguspermana/quota-race/internal/config"
)

type Result struct {
	Status int
	Err    string
	Body   string
}

type Attempt struct {
	Index      int
	Results    []Result
	ByStatus   map[int]int
	OKCount    int
	ErrorCount int
	Elapsed    time.Duration
	FollowBody string
	FollowStat int
}

func Run(ctx context.Context, cfg config.File, client *http.Client) ([]Attempt, error) {
	if client == nil {
		client = &http.Client{Timeout: time.Duration(cfg.TimeoutMS) * time.Millisecond}
	}
	out := make([]Attempt, 0, cfg.Attempts)
	for i := 0; i < cfg.Attempts; i++ {
		if cfg.Reset != nil && cfg.Reset.URL != "" {
			_ = doReq(ctx, client, *cfg.Reset)
		}
		a, err := oneAttempt(ctx, cfg, client, i)
		if err != nil {
			return out, err
		}
		out = append(out, a)
	}
	return out, nil
}

func oneAttempt(ctx context.Context, cfg config.File, client *http.Client, index int) (Attempt, error) {
	start := time.Now()
	n := cfg.Concurrency
	results := make([]Result, n)
	var startGate sync.WaitGroup
	startGate.Add(1)
	var done sync.WaitGroup
	done.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer done.Done()
			startGate.Wait()
			results[i] = doReq(ctx, client, cfg.Request)
		}()
	}
	startGate.Done()
	done.Wait()

	a := Attempt{
		Index:    index,
		Results:  results,
		ByStatus: map[int]int{},
		Elapsed:  time.Since(start),
	}
	for _, r := range results {
		if r.Err != "" {
			a.ErrorCount++
			continue
		}
		a.ByStatus[r.Status]++
		if r.Status == cfg.Invariant.OKStatus {
			a.OKCount++
		}
	}
	if cfg.FollowUp != nil && cfg.FollowUp.URL != "" {
		fr := doReq(ctx, client, *cfg.FollowUp)
		a.FollowStat = fr.Status
		a.FollowBody = fr.Body
		if fr.Err != "" && a.FollowBody == "" {
			a.FollowBody = fr.Err
		}
	}
	return a, nil
}

func doReq(ctx context.Context, client *http.Client, req config.Request) Result {
	var body io.Reader
	if req.Body != "" {
		body = bytes.NewBufferString(req.Body)
	}
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, body)
	if err != nil {
		return Result{Err: err.Error()}
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	if req.Body != "" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return Result{Err: err.Error()}
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	return Result{Status: resp.StatusCode, Body: string(b)}
}
