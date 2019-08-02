package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os/exec"
)

// Not matching paths will not be served
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.serverError(w, errors.New(r.URL.Path+" not implemented"))
		return
	}
}

// Check if Docker daemon is running.
func (app *application) commandDockerInfo(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("/bin/sh", "-c", "docker info")
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Printf("Command finished with error: %v", stderr.String())
		app.serverError(w, errors.New(stderr.String()))
		return
	}

	app.serverOk(w)
}
