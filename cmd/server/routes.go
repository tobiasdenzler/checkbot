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

	return mux
}
