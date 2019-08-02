package main

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

// Check if Docker daemon is running.
func (app *application) checkDockerInfo(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("docker info", w, r)
}

// Check if Docker pull is working.
func (app *application) checkDockerPull(w http.ResponseWriter, r *http.Request) {
	app.runShellCommand("docker pull hello-world", w, r)
}

// Check if OpenShift master API is working.
func (app *application) checkMasterAPI(w http.ResponseWriter, r *http.Request) {
	app.runCurlCommand("https://192.168.42.17:8443/healthz", w, r)
}

// Run a shell command and return a http result.
func (app *application) runShellCommand(command string, w http.ResponseWriter, r *http.Request) {

	app.infoLog.Printf("Execute shell command: %s", command)

	cmd := exec.Command("/bin/sh", "-c", command)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		app.infoLog.Printf("Command finished with error: %v", stderr.String())
		app.serverError(w, errors.New(stderr.String()))
		return
	}

	app.infoLog.Println(out.String())
	app.serverTrue(w)
}

// Run a shell command and return a http result.
func (app *application) runCurlCommand(url string, w http.ResponseWriter, r *http.Request) {

	app.infoLog.Printf("Execute curl command: %s", url)

	cmd := exec.Command("curl", "-k", "-s", "-o", "/dev/null", "-w", "'%{http_code}'", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		app.infoLog.Printf("Command finished with error: %v", err.Error())
		app.serverError(w, err)
		return
	}

	// Get status code from curl response
	statuscode, err := strconv.Atoi(strings.Trim(out.String(), "'"))
	app.infoLog.Printf("Result from command: %d", statuscode)

	if statuscode >= 400 {
		app.serverFalse(w)
		return
	}
	app.serverTrue(w)
}
