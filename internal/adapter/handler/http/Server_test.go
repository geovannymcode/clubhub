package http

import (
	"github.com/go-chi/chi/v5"
	"net/http/httptest"
	"testing"
)

func TestInvalidAddressRejectedBeforeIO(t *testing.T) {
	handler := &Domain{} // nil repo: rejection must happen before any IO.
	router := chi.NewRouter()
	router.Get("/address={address}", handler.GetByAddress)
	for _, address := range []string{"127.0.0.1", "localhost", "example.com:443", "example.com%3Fhost=localhost"} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", "/address="+address, nil))
		if w.Code != 400 {
			t.Errorf("%s status %d", address, w.Code)
		}
	}
}
