package main

import (
	"net/http"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	app.serverTrue(w)
}

// Check if all projects have quota.
func (app *application) checkOpenshiftProjectQuota(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("./scripts/compare-quota.sh", w, r)
}
