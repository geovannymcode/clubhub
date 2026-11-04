package security

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

const testToken = "0123456789abcdef0123456789abcdef"

func request(method, token string) *http.Request {
	r := httptest.NewRequest(method, "http://example.com/clubhub", nil)
	if token != "" {
		r.Header.Set("Authorization", token)
	}
	return r
}

func TestProtectAuthenticationAndHeaders(t *testing.T) {
	calls := 0
	protected := Protect(testToken, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.WriteHeader(204) }))
	for _, auth := range []string{"", "Bearer wrong", "Basic " + testToken, testToken, "Bearer " + testToken + "extra"} {
		w := httptest.NewRecorder()
		protected.ServeHTTP(w, request("GET", auth))
		if w.Code != 401 {
			t.Errorf("auth %q status %d", auth, w.Code)
		}
		if strings.Contains(w.Body.String(), testToken) {
			t.Fatal("credential leaked")
		}
	}
	if calls != 0 {
		t.Fatal("unauthorized request reached handler")
	}
	w := httptest.NewRecorder()
	r := request("GET", "Bearer "+testToken)
	r.Header.Set("Origin", "https://evil.example")
	protected.ServeHTTP(w, r)
	if w.Code != 204 || calls != 1 {
		t.Fatalf("valid auth failed: %d calls %d", w.Code, calls)
	}
	for _, key := range []string{"X-Content-Type-Options", "X-Frame-Options", "Content-Security-Policy", "Cache-Control"} {
		if w.Header().Get(key) == "" {
			t.Errorf("missing header %s", key)
		}
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("cross origin access enabled")
	}
	w = httptest.NewRecorder()
	Protect("", protected).ServeHTTP(w, r)
	if w.Code != 401 {
		t.Fatal("empty configured token accepted")
	}
}

func TestProtectRequestBudget(t *testing.T) {
	protected := Protect(testToken, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	for i := 0; i < 21; i++ {
		w := httptest.NewRecorder()
		protected.ServeHTTP(w, request("GET", "Bearer "+testToken))
		if i < 20 && w.Code != 204 {
			t.Fatalf("request %d status %d", i, w.Code)
		}
		if i == 20 && (w.Code != 429 || w.Header().Get("Retry-After") == "") {
			t.Fatal("rate limit absent")
		}
	}
}

func TestProtectMethodsAndBodies(t *testing.T) {
	protected := Protect(testToken, http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Fatal("invalid request reached handler") }))
	for _, method := range []string{"POST", "PUT", "DELETE", "OPTIONS"} {
		w := httptest.NewRecorder()
		protected.ServeHTTP(w, request(method, "Bearer "+testToken))
		if w.Code != 405 {
			t.Fatalf("%s status %d", method, w.Code)
		}
	}
	r := request("GET", "Bearer "+testToken)
	r.ContentLength = 10
	w := httptest.NewRecorder()
	protected.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatal("GET body accepted")
	}
}

func TestProtectConcurrency(t *testing.T) {
	entered := make(chan struct{}, 8)
	release := make(chan struct{})
	protected := Protect(testToken, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { entered <- struct{}{}; <-release }))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			protected.ServeHTTP(httptest.NewRecorder(), request("GET", "Bearer "+testToken))
		}()
	}
	for i := 0; i < 8; i++ {
		<-entered
	}
	w := httptest.NewRecorder()
	protected.ServeHTTP(w, request("GET", "Bearer "+testToken))
	close(release)
	wg.Wait()
	if w.Code != 429 {
		t.Fatal("concurrency limit absent")
	}
}
