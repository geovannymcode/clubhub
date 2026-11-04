package security

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"testing"
)

func TestDomain(t *testing.T) {
	for input, want := range map[string]string{"Example.COM": "example.com", "https://example.com/": "example.com", "example.com.": "example.com", "xn--bcher-kva.de": "xn--bcher-kva.de"} {
		got, err := Domain(input)
		if err != nil || got != want {
			t.Errorf("Domain(%q) = %q, %v", input, got, err)
		}
	}
	for _, input := range []string{"", "localhost", "127.0.0.1", "[::1]", "https://user@example.com", "https://example.com:443", "example.com:80", "example.com/path", "example.com?host=localhost", "example.com#test", "https://example.com/?x=y", "https://example.com/#fragment", "https://example.com/../admin", "evil..com", "-evil.com", "evil-.com", "example.com\r\nHost: localhost", strings.Repeat("a", 64) + ".com"} {
		if _, err := Domain(input); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
}

func TestPublicIP(t *testing.T) {
	for _, value := range []string{"127.0.0.1", "10.1.2.3", "172.16.0.1", "192.168.1.1", "169.254.169.254", "100.100.100.200", "0.0.0.0", "224.0.0.1", "255.255.255.255", "192.0.0.8", "192.0.2.1", "198.18.0.1", "198.51.100.1", "203.0.113.1", "::1", "::", "fc00::1", "fe80::1", "::ffff:127.0.0.1", "64:ff9b::a00:1", "2001:db8::1", "2002:7f00:1::", "3fff::1"} {
		if PublicIP(netip.MustParseAddr(value)) {
			t.Errorf("allowed %s", value)
		}
	}
	for _, value := range []string{"1.1.1.1", "8.8.8.8", "2606:4700:4700::1111", "::ffff:8.8.8.8"} {
		if !PublicIP(netip.MustParseAddr(value)) {
			t.Errorf("blocked %s", value)
		}
	}
}

func TestPublicDialPinsValidatedDNS(t *testing.T) {
	lookups := 0
	lookup := func(ctx context.Context, host string) ([]netip.Addr, error) {
		lookups++
		if host != "example.com" {
			t.Fatalf("unexpected DNS host %q", host)
		}
		return []netip.Addr{netip.MustParseAddr("8.8.8.8")}, nil
	}
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr != "8.8.8.8:443" {
			t.Fatalf("dial was not pinned: %s", addr)
		}
		return nil, errors.New("test dial failure")
	}
	publicDial(context.Background(), "tcp", "example.com:443", lookup, dial)
	if lookups != 1 {
		t.Fatalf("resolved %d times", lookups)
	}
}

func TestPublicDialRejectsMixedAnswers(t *testing.T) {
	lookup := func(context.Context, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("8.8.8.8"), netip.MustParseAddr("127.0.0.1")}, nil
	}
	dialed := false
	_, err := publicDial(context.Background(), "tcp", "example.com:443", lookup, func(context.Context, string, string) (net.Conn, error) { dialed = true; return nil, nil })
	if err == nil || dialed {
		t.Fatal("dialed a mixed public/private DNS answer")
	}
}

func TestPrivateDialBlocked(t *testing.T) {
	for _, target := range []string{"127.0.0.1:443", "169.254.169.254:443", "[::1]:443"} {
		if conn, err := DialPublic(context.Background(), "tcp", target); err == nil {
			conn.Close()
			t.Fatalf("connected to %s", target)
		}
	}
}

func TestRedirectPolicy(t *testing.T) {
	client := NewHTTPClient()
	for _, raw := range []string{"http://example.com/", "https://user:pass@example.com/", "https://example.com:8443/"} {
		u, _ := url.Parse(raw)
		if err := client.CheckRedirect(&http.Request{URL: u}, nil); err == nil {
			t.Errorf("accepted redirect %s", raw)
		}
	}
	u, _ := url.Parse("https://example.com/")
	if err := client.CheckRedirect(&http.Request{URL: u}, make([]*http.Request, 3)); err == nil {
		t.Fatal("redirect limit absent")
	}
	tr := client.Transport.(*http.Transport)
	if tr.Proxy != nil {
		t.Fatal("environment proxy can bypass destination checks")
	}
}

func TestReadLimited(t *testing.T) {
	if _, err := ReadLimited(strings.NewReader("12345"), 4); err == nil {
		t.Fatal("accepted oversized body")
	}
	if b, err := ReadLimited(strings.NewReader("1234"), 4); err != nil || string(b) != "1234" {
		t.Fatalf("boundary: %q %v", b, err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestFetchStatusAndLimits(t *testing.T) {
	original := outbound
	defer func() { outbound = original }()
	for _, tc := range []struct {
		status int
		body   string
	}{{500, "upstream secret"}, {200, strings.Repeat("x", (2<<20)+1)}} {
		outbound = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
		})}
		if _, err := Fetch(context.Background(), "https://example.com/"); err == nil {
			t.Fatalf("accepted status %d / length %d", tc.status, len(tc.body))
		}
	}
}
