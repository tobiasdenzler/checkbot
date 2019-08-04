package main

import "net/http"

// Define the routes to serve
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Daemon checks are running directly on the node
	mux.HandleFunc("/daemon/docker/info", app.checkDockerInfo)
	mux.HandleFunc("/daemon/docker/pull", app.checkDockerPull)

	// Server checks are running in a container on the cluster
	mux.HandleFunc("/server/openshift/masterapi", app.checkOpenshiftMasterAPI)
	mux.HandleFunc("/server/openshift/console", app.checkOpenshiftConsole)
	mux.HandleFunc("/server/openshift/router", app.checkOpenshiftRouter)
	mux.HandleFunc("/server/openshift/nodes", app.checkOpenshiftNodes)
	mux.HandleFunc("/server/openshift/projectquota", app.checkOpenshiftProjectQuota)

	return mux
}
