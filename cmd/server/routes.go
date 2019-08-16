package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define the routes to serve
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Health endpoint
	mux.HandleFunc("/health", health)

	return mux
}

// Health check of server
func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
