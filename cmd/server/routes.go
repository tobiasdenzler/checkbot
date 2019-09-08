package main

import (
	"net/http"

	"github.com/goji/httpauth"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define the routes to serve
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/sandbox", app.sandbox)
	mux.HandleFunc("/health", app.health)

	mux.Handle("/reload", httpauth.SimpleBasicAuth("admin", "admin")(http.HandlerFunc(app.reload)))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
