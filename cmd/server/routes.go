package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Define the routes to serve
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Reload endpoint
	mux.HandleFunc("/reload", app.reload)

	// Health endpoint
	mux.HandleFunc("/health", health)

	return mux
}

// Reload all chekcs
func (app *application) reload(w http.ResponseWriter, r *http.Request) {
	log.Info("Reloading checks..")

	// Stop all checks
	app.stopChecks()
	// Rebuild all checks
	app.buildMetrics()
	// Start all checks
	app.startChecks()
}

// Health check of server
func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
