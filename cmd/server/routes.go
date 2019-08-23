package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define the routes to serve
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/reload", app.reload)
	mux.HandleFunc("/sandbox", app.sandbox)
	mux.HandleFunc("/health", app.health)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
