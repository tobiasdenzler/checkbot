package main

import (
	"net/http"

	"github.com/goji/httpauth"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define the routes to serve
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)

	// Metrics endpoint for Prometheus
	mux.Handle("/metrics", promhttp.Handler())

	// Sandbox
	if app.config.Sandbox {
		mux.Handle("/sandbox", httpauth.SimpleBasicAuth("admin", app.managementPwd)(http.HandlerFunc(app.sandbox)))
	}

	// Health endpoint
	mux.HandleFunc("/health", app.health)

	// Reload scripts endpoint
	mux.Handle("/reload", httpauth.SimpleBasicAuth("admin", app.managementPwd)(http.HandlerFunc(app.reload)))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return app.logRequest(secureHeaders(mux))
}
