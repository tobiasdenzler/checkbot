package main

import (
	"net/http"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	app.serverTrue(w)
}

// Check if all projects have quota.
func (app *application) checkProjectQuota(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("./scripts/projectHasQuota.sh", w, r)
}

// Check if Daemonset is running.
func (app *application) checkDaemonsetRunning(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("./scripts/daemonsetIsRunning.sh", w, r)
}
