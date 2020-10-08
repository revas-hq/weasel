package http

import (
	"net/http"
	"strings"
)

var allowMethods = "GET, HEAD, OPTIONS"

func NewCORSMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hs := w.Header()
			hs.Set("allow", allowMethods)
			hs.Set("Access-Control-Allow-Origin", "*")
			if r.Method == "OPTIONS" {
				hs.Set("Access-Control-Allow-Methods", allowMethods)
				hs.Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
				hs.Set("Access-Control-Expose-Headers", "Location, Etag, Content-Disposition")
				// hs.Set("access-control-max-age", "86400")
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func NewCheckMethodMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(allowMethods, r.Method) {
				http.Error(w, "", http.StatusMethodNotAllowed)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

const stsValue = "max-age=10886400; includeSubDomains; preload;"

func NewForceTLSMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Forwarded-Proto") == "" {
				http.Error(w, "", http.StatusTeapot)
				return
			}
			if r.Header.Get("X-Forwarded-Proto") == "http" {
				u := "https://" + r.Host + r.URL.Path
				if r.URL.RawQuery != "" {
					u += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, u, http.StatusMovedPermanently)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func NewForceSTSMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Forwarded-Proto") == "" {
				http.Error(w, "", http.StatusTeapot)
				return
			}
			if r.Header.Get("X-Forwarded-Proto") == "https" {
				w.Header().Set("Strict-Transport-Security", stsValue)
			}
			h.ServeHTTP(w, r)
		})
	}
}
