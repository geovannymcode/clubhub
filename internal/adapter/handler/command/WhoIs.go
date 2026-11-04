package command

import (
	"context"
	"fmt"
	"github.com/Geovanny0401/clubhub/internal/security"
	"net/netip"
	"strings"
	"time"
)

// WHOIS enrichment is best effort. Query only a fixed registry, never referrals,
// with a bounded response and a deadline inherited from the HTTP request.
func RunWhoIs(ctx context.Context, ipAddr string) map[string][]string {
	result := make(map[string][]string)
	ip, err := netip.ParseAddr(ipAddr)
	if err != nil || !security.PublicIP(ip) {
		return result
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := security.DialPublic(ctx, "tcp", "whois.arin.net:43")
	if err != nil {
		return result
	}
	defer conn.Close()
	deadline, _ := ctx.Deadline()
	if err := conn.SetDeadline(deadline); err != nil {
		return result
	}
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()
	if _, err := fmt.Fprintf(conn, "n + %s\r\n", ip.String()); err != nil {
		return result
	}
	body, err := security.ReadLimited(conn, 256<<10)
	if err != nil {
		return result
	}
	for _, line := range strings.Split(string(body), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if ok && (key == "Country" || key == "OrgName") {
			result[key] = []string{strings.TrimSpace(value)}
		}
	}
	return result
}
