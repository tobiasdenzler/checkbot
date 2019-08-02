package main

import "net/http"

// Define the routes to serve
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/check/dockerinfo", app.checkDockerInfo)
	mux.HandleFunc("/check/dockerpull", app.checkDockerPull)
	mux.HandleFunc("/check/masterapi", app.checkMasterAPI)

	return mux
}
