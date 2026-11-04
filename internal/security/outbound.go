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
	"time"
)

// In addition to private and local addresses, exclude special-use networks that
// may reach shared infrastructure, metadata services or non-public hosts.
var specialNetworks = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"), netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("192.0.0.0/24"), netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"), netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"), netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("240.0.0.0/4"), netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:db8::/32"), netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("3fff::/20"),
}

func PublicIP(ip netip.Addr) bool {
	ip = ip.Unmap()
	if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip.Is6() && !netip.MustParsePrefix("2000::/3").Contains(ip) {
		return false
	}
	for _, network := range specialNetworks {
		if network.Contains(ip) {
			return false
		}
	}
	return true
}

// Domain accepts DNS names only: no credentials, paths, IP literals or ports.
func Domain(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://") {
		u, err := url.Parse(value)
		if err != nil || u.User != nil || u.Port() != "" || (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
			return "", errors.New("invalid domain")
		}
		value = u.Hostname()
	}
	value = strings.TrimSuffix(value, ".")
	if len(value) > 100 || !strings.Contains(value, ".") {
		return "", errors.New("invalid domain")
	}
	if _, err := netip.ParseAddr(value); err == nil {
		return "", errors.New("IP addresses are not allowed")
	}
	for _, label := range strings.Split(value, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return "", errors.New("invalid domain")
		}
		for _, ch := range label {
			if !(ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '-') {
				return "", errors.New("invalid domain")
			}
		}
	}
	return value, nil
}

type lookupFunc func(context.Context, string) ([]netip.Addr, error)
type dialFunc func(context.Context, string, string) (net.Conn, error)

// Resolve once and dial the validated numeric address to prevent DNS rebinding.
func publicDial(ctx context.Context, network, address string, lookup lookupFunc, dial dialFunc) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	ips, err := lookup(ctx, host)
	if err != nil || len(ips) == 0 {
		return nil, errors.New("cannot resolve public host")
	}
	for _, ip := range ips {
		if !PublicIP(ip) {
			return nil, errors.New("non-public destination blocked")
		}
	}
	for _, ip := range ips {
		conn, dialErr := dial(ctx, network, net.JoinHostPort(ip.String(), port))
		if dialErr == nil {
			return conn, nil
		}
		err = dialErr
	}
	return nil, err
}

func DialPublic(ctx context.Context, network, address string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	lookup := func(ctx context.Context, host string) ([]netip.Addr, error) {
		return net.DefaultResolver.LookupNetIP(ctx, "ip", host)
	}
	return publicDial(ctx, network, address, lookup, dialer.DialContext)
}

func validateURL(u *url.URL) error {
	if u.Scheme != "https" || u.User != nil || (u.Port() != "" && u.Port() != "443") || u.Hostname() == "" {
		return errors.New("only public HTTPS destinations on port 443 are allowed")
	}
	return nil
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext:            DialPublic,
			TLSHandshakeTimeout:    3 * time.Second,
			ResponseHeaderTimeout:  5 * time.Second,
			MaxResponseHeaderBytes: 32 << 10,
			MaxIdleConns:           16, MaxIdleConnsPerHost: 2,
			IdleConnTimeout: 30 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return errors.New("too many redirects")
			}
			return validateURL(req.URL)
		},
	}
}

var outbound = NewHTTPClient()

func Fetch(ctx context.Context, address string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	if err := validateURL(req.URL); err != nil {
		return nil, err
	}
	resp, err := outbound.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("upstream returned an unsuccessful status")
	}
	return ReadLimited(resp.Body, 2<<20)
}

func ReadLimited(r io.Reader, limit int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, errors.New("upstream response too large")
	}
	return data, nil
}
