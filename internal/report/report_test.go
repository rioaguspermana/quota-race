package report

import (
	"testing"

	"github.com/rioaguspermana/quota-race/internal/burst"
	"github.com/rioaguspermana/quota-race/internal/config"
)

func TestSummarizeFailWhenTooManyGrants(t *testing.T) {
	cfg := config.File{Invariant: config.Invariant{OKStatus: 200, MaxOK: 10}}
	attempts := []burst.Attempt{{
		OKCount:  40,
		ByStatus: map[int]int{200: 40},
	}}
	s := Summarize(cfg, attempts)
	if s.Pass {
		t.Fatal("expected fail")
	}
}

func TestSummarizePassAtLimit(t *testing.T) {
	cfg := config.File{Invariant: config.Invariant{OKStatus: 200, MaxOK: 10}}
	attempts := []burst.Attempt{{
		OKCount:  10,
		ByStatus: map[int]int{200: 10, 409: 30},
	}}
	s := Summarize(cfg, attempts)
	if !s.Pass {
		t.Fatal("expected pass")
	}
}
