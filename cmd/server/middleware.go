package main

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if(!strings.Contains(r.URL.RequestURI(), "/health") && 
			!strings.Contains(r.URL.RequestURI(), "/metrics")) {
			log.Tracef("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		}
		next.ServeHTTP(w, r)
	})
}
