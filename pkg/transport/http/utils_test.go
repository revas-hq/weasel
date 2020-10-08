package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	weasel_http "github.com/revas-hq/weasel/pkg/transport/http"
)

var NopHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

func GetForwardedProto(url string) string {
	if strings.HasPrefix(url, "https") {
		return "https"
	}
	if strings.HasPrefix(url, "http") {
		return "http"
	}
	return ""
}

func TestForceTLSMiddleware(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		url      string
		code     int
		location string
	}{
		{"http://example.com", http.StatusMovedPermanently, "https://example.com"},
		{"http://example.com/some/path.html", http.StatusMovedPermanently, "https://example.com/some/path.html"},
		{"http://example.com/?query=value", http.StatusMovedPermanently, "https://example.com/?query=value"},
		{"http://example.com/some/path.html?query=value", http.StatusMovedPermanently, "https://example.com/some/path.html?query=value"},
		{"https://example.com", http.StatusOK, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.url, func(t *testing.T) {
			t.Parallel()
			r, _ := http.NewRequest("GET", tt.url, nil)
			r.Header.Set("X-Forwarded-Proto", GetForwardedProto(tt.url))
			w := httptest.NewRecorder()

			h := weasel_http.NewForceTLSMiddleware()(NopHandler)
			h.ServeHTTP(w, r)

			if w.Code != tt.code {
				t.Errorf("w.Code = %d; want %d", w.Code, tt.code)
			}
			if v := w.Header().Get("Location"); v != tt.location {
				t.Errorf("Location: %s; want %s", v, tt.location)
			}
		})
	}
}

func TestForceSTSMiddleware(t *testing.T) {
	h := weasel_http.NewForceTLSMiddleware()(NopHandler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("origin", "https://example.com")

	h.ServeHTTP(w, r)
}
