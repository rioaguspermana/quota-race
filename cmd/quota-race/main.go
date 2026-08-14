package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/rioaguspermana/quota-race/internal/burst"
	"github.com/rioaguspermana/quota-race/internal/config"
	"github.com/rioaguspermana/quota-race/internal/report"
	"github.com/rioaguspermana/quota-race/internal/safety"
)

func main() {
	os.Exit(run())
}

func run() int {
	fs := flag.NewFlagSet("quota-race", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	configPath := fs.String("c", "", "path to YAML config")
	jsonOut := fs.Bool("json", false, "print JSON report")
	own := fs.Bool("i-own-this-api", false, "required for non-loopback URLs; you own the target or have written permission")
	raw := os.Args[1:]
	if len(raw) == 0 || raw[0] == "help" || raw[0] == "-h" || raw[0] == "--help" {
		usage(fs)
		return 0
	}
	if raw[0] != "run" {
		fmt.Fprintf(os.Stderr, "unknown command %q\n", raw[0])
		usage(fs)
		return 2
	}
	if err := fs.Parse(raw[1:]); err != nil {
		return 2
	}
	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "missing -c config.yaml")
		return 2
	}

	safety.PrintBanner(os.Stderr)

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		return 2
	}
	ownAPI := *own || cfg.IOwnThisAPI
	if err := safety.CheckTarget(cfg.Request.URL, ownAPI); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 2
	}
	if cfg.FollowUp != nil && cfg.FollowUp.URL != "" {
		if err := safety.CheckTarget(cfg.FollowUp.URL, ownAPI); err != nil {
			fmt.Fprintf(os.Stderr, "follow_up: %v\n", err)
			return 2
		}
	}
	if cfg.Reset != nil && cfg.Reset.URL != "" {
		if err := safety.CheckTarget(cfg.Reset.URL, ownAPI); err != nil {
			fmt.Fprintf(os.Stderr, "reset: %v\n", err)
			return 2
		}
	}

	attempts, err := burst.Run(context.Background(), cfg, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "run: %v\n", err)
		return 1
	}
	sum := report.Summarize(cfg, attempts)
	if *jsonOut {
		if err := report.WriteJSON(os.Stdout, sum); err != nil {
			fmt.Fprintf(os.Stderr, "json: %v\n", err)
			return 1
		}
	} else {
		report.WriteText(os.Stdout, sum)
	}
	if !sum.Pass {
		return 1
	}
	return 0
}

func usage(fs *flag.FlagSet) {
	fmt.Fprintf(os.Stderr, "quota-race — check a quota/stock/seat invariant under concurrent HTTP requests.\n\n")
	safety.PrintBanner(os.Stderr)
	fmt.Fprintf(os.Stderr, "Usage:\n  quota-race run -c examples/racey.yaml\n  quota-race run -c examples/fixed.yaml --json\n\n")
	fmt.Fprintf(os.Stderr, "Not a k6 replacement. Not for third-party production.\n\n")
	fs.PrintDefaults()
}
