package main

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

// Run a shell command and return a http result.
func (app *application) runShellCommand(command string, w http.ResponseWriter, r *http.Request) {

	app.infoLog.Printf("Execute shell command: %s", command)

	cmd := exec.Command("/bin/sh", command)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		// Check failed with defined message
		if out.String() != "" {
			app.infoLog.Printf("Check %s failed with output: %v", command, out.String())
			app.serverError(w, errors.New(out.String()))
			return
		}

		// Execution failed
		app.infoLog.Printf("Command %s finished with error: %v", command, stderr.String())
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
		app.serverFalse(w)
		return
	}

	// Get status code from curl response
	statuscode, err := strconv.Atoi(strings.Trim(out.String(), "'"))
	app.infoLog.Printf("Result from command: %d", statuscode)

	if statuscode >= 400 || statuscode == 0 {
		app.serverFalse(w)
		return
	}
	app.serverTrue(w)
}

// Returns a false check result
func (app *application) serverError(w http.ResponseWriter, err error) {
	app.errorLog.Output(2, err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Returns a 200 and true check result
func (app *application) serverTrue(w http.ResponseWriter) {
	w.Write([]byte("true"))
}

// Returns a 404 and false check result
func (app *application) serverFalse(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("false"))
}
