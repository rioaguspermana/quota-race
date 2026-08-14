package safety

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

const Banner = `quota-race checks a business invariant under concurrent requests.
Use it only against APIs you own, local servers, or staging you have
explicit permission to test. It is not a load tester and not a way to
probe third-party production.

Demo servers in this repo (examples/counter) are the intended first target.
`

const MaxConcurrency = 256

func PrintBanner(w io.Writer) {
	fmt.Fprint(w, Banner)
	fmt.Fprintln(w)
}

func HostIsLoopback(raw string) (bool, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return false, fmt.Errorf("url: %w", err)
	}
	host := u.Hostname()
	if host == "" {
		return false, fmt.Errorf("url missing host")
	}
	if strings.EqualFold(host, "localhost") {
		return true, nil
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return true, nil
	}
	return false, nil
}

func CheckTarget(rawURL string, ownAPI bool) error {
	loopback, err := HostIsLoopback(rawURL)
	if err != nil {
		return err
	}
	if loopback {
		return nil
	}
	if !ownAPI {
		return fmt.Errorf("refusing non-loopback URL %q: pass --i-own-this-api (and set i_own_this_api in the config) — only APIs you own or have written permission to test", rawURL)
	}
	return nil
}
