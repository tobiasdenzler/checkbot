package main

import (
	"bytes"
	"errors"
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
func (app *application) checkDockerInfo(w http.ResponseWriter, r *http.Request) {

	out, stderr, err := app.runShellCommand("docker info")

	if err != nil {
		app.infoLog.Printf("Command failed: %v", stderr)
		app.serverError(w, errors.New(stderr))
		return
	}

	app.infoLog.Printf("Command succeeded: %v", out)
	app.serverOk(w)
}

// Check if Docker pull is working.
func (app *application) checkDockerPull(w http.ResponseWriter, r *http.Request) {

	out, stderr, err := app.runShellCommand("docker pull hello-world")

	if err != nil {
		app.infoLog.Printf("Command finished with error: %v", stderr)
		app.serverError(w, errors.New(stderr))
		return
	}

	app.infoLog.Println(out)
	app.serverOk(w)
}

func (app *application) runShellCommand(command string) (string, string, error) {
	app.infoLog.Printf("Execute shell command: %s", command)
	cmd := exec.Command("/bin/sh", "-c", command)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}
