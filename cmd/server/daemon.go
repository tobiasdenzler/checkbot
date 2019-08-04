package main

import (
	"net/http"
)

// Check if Docker daemon is running.
func (app *application) checkDockerInfo(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("docker info", w, r)
}

// Check if Docker pull is working.
func (app *application) checkDockerPull(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("docker pull hello-world", w, r)
}
