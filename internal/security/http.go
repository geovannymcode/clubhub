package security

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Protect applies a single shared API credential, a global request budget and
// a hard concurrency bound before any database or external network work.
func Protect(token string, next http.Handler) http.Handler {
	expected := sha256.Sum256([]byte(token))
	var mu sync.Mutex
	budget, last := 20.0, time.Now()
	slots := make(chan struct{}, 8)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		if r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		}
		auth := r.Header.Get("Authorization")
		provided := sha256.Sum256([]byte(strings.TrimPrefix(auth, "Bearer ")))
		if len(token) < 32 || !strings.HasPrefix(auth, "Bearer ") || subtle.ConstantTimeCompare(provided[:], expected[:]) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		mu.Lock()
		now := time.Now()
		budget = min(20, budget+now.Sub(last).Seconds())
		last = now
		allowed := budget >= 1
		if allowed {
			budget--
		}
		mu.Unlock()
		if !allowed {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		select {
		case slots <- struct{}{}:
			defer func() { <-slots }()
		default:
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.ContentLength != 0 || len(r.TransferEncoding) != 0 {
			http.Error(w, "Request body is not allowed", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
