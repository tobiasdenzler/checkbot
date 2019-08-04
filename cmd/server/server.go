package main

import (
	"net/http"
)

// Check if OpenShift master API is working.
func (app *application) checkOpenshiftMasterAPI(w http.ResponseWriter, r *http.Request) {
	app.runCurlCommand("https://"+app.openshiftHost+":8443/healthz", w, r)
}

// Check if OpenShift console is working.
func (app *application) checkOpenshiftConsole(w http.ResponseWriter, r *http.Request) {
	app.runCurlCommand("http://"+app.openshiftHost+":8443/console/", w, r)
}

// Check if OpenShift router is working.
func (app *application) checkOpenshiftRouter(w http.ResponseWriter, r *http.Request) {
	app.runCurlCommand("http://"+app.openshiftHost+":1936/healthz", w, r)
}

// Check if oc get nodes is working.
func (app *application) checkOpenshiftNodes(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("oc get nodes -o wide", w, r)
}

// Check if all projects have quota.
func (app *application) checkOpenshiftProjectQuota(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("./compare-quota.sh", w, r)
}
