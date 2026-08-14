package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/rioaguspermana/quota-race/internal/safety"
	"gopkg.in/yaml.v3"
)

type File struct {
	IOwnThisAPI bool      `yaml:"i_own_this_api"`
	Request     Request   `yaml:"request"`
	Concurrency int       `yaml:"concurrency"`
	Attempts    int       `yaml:"attempts"`
	TimeoutMS   int       `yaml:"timeout_ms"`
	Invariant   Invariant `yaml:"invariant"`
	FollowUp    *Request  `yaml:"follow_up"`
	Reset       *Request  `yaml:"reset"`
}

type Request struct {
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

type Invariant struct {
	// Status counted as a "grant" (success that consumes quota).
	OKStatus int `yaml:"ok_status"`
	// MaxOK is the business limit: at most this many OKStatus responses per attempt.
	MaxOK int `yaml:"max_ok"`
}

func Load(path string) (File, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return File{}, err
	}
	var cfg File
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return File{}, fmt.Errorf("yaml: %w", err)
	}
	if err := cfg.normalize(); err != nil {
		return File{}, err
	}
	return cfg, nil
}

func (c *File) normalize() error {
	if strings.TrimSpace(c.Request.URL) == "" {
		return fmt.Errorf("request.url is required")
	}
	if c.Request.Method == "" {
		c.Request.Method = "POST"
	} else {
		c.Request.Method = strings.ToUpper(c.Request.Method)
	}
	if c.Concurrency <= 0 {
		c.Concurrency = 10
	}
	if c.Concurrency > safety.MaxConcurrency {
		return fmt.Errorf("concurrency %d exceeds hard cap %d (this is an invariant checker, not a load weapon)", c.Concurrency, safety.MaxConcurrency)
	}
	if c.Attempts <= 0 {
		c.Attempts = 1
	}
	if c.Attempts > 20 {
		return fmt.Errorf("attempts %d exceeds cap 20", c.Attempts)
	}
	if c.TimeoutMS <= 0 {
		c.TimeoutMS = 5000
	}
	if c.Invariant.OKStatus == 0 {
		c.Invariant.OKStatus = 200
	}
	if c.Invariant.MaxOK < 0 {
		return fmt.Errorf("invariant.max_ok must be >= 0")
	}
	if c.FollowUp != nil {
		if c.FollowUp.Method == "" {
			c.FollowUp.Method = "GET"
		} else {
			c.FollowUp.Method = strings.ToUpper(c.FollowUp.Method)
		}
	}
	if c.Reset != nil {
		if c.Reset.Method == "" {
			c.Reset.Method = "POST"
		} else {
			c.Reset.Method = strings.ToUpper(c.Reset.Method)
		}
	}
	return nil
}
