package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Overview of all checks
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.render(w, r, "checks.page.tmpl", &templateData{
		Checklist:      app.checkList,
		SandboxEnabled: app.enableSandbox,
	})
}

// Render sandbox form for debugging
func (app *application) sandbox(w http.ResponseWriter, r *http.Request) {

	var sandbox Sandbox

	// POST from form
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		// Load existing script to sandbox
		if r.PostForm.Get("load") != "none" {
			check := app.checkList[r.PostForm.Get("load")]
			r.PostForm.Set("sandbox", app.loadSandbox(*check))
		}

		// Execute script in sandbox
		sandbox = *app.runSandbox(r.PostForm.Get("sandbox"))
	}

	app.render(w, r, "sandbox.page.tmpl", &templateData{
		Checklist:      app.checkList,
		Sandbox:        sandbox,
		SandboxEnabled: app.enableSandbox,
	})
}

// Reload all chekcs
func (app *application) reload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Info("Reloading checks..")
		// Stop all checks
		app.stopChecks()
		// Rebuild all checks
		app.buildMetrics()
		// Start all checks
		app.startChecks()
	} else {
		http.NotFound(w, r)
	}
}

// Health check of server
func (app *application) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
