package safety

import "testing"

func TestHostIsLoopback(t *testing.T) {
	ok, err := HostIsLoopback("http://127.0.0.1:8090/take")
	if err != nil || !ok {
		t.Fatalf("127.0.0.1: %v %v", ok, err)
	}
	ok, err = HostIsLoopback("http://localhost:8090/take")
	if err != nil || !ok {
		t.Fatalf("localhost: %v %v", ok, err)
	}
	ok, err = HostIsLoopback("https://example.com/quota")
	if err != nil || ok {
		t.Fatalf("example.com should not be loopback: %v %v", ok, err)
	}
}

func TestCheckTarget(t *testing.T) {
	if err := CheckTarget("http://127.0.0.1:8090/x", false); err != nil {
		t.Fatal(err)
	}
	if err := CheckTarget("https://api.example.com/x", false); err == nil {
		t.Fatal("expected refuse")
	}
	if err := CheckTarget("https://api.example.com/x", true); err != nil {
		t.Fatal(err)
	}
}
