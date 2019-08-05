package main

import "net/http"

// Define the routes to serve
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Internal health check
	mux.HandleFunc("/health", app.health)

	// Server checks
	mux.HandleFunc("/check/projectHasQuota", app.checkProjectQuota)
	mux.HandleFunc("/check/daemonsetIsRunning", app.checkDaemonsetRunning)

	return mux
}
