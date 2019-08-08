package main

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	app.serverTrue(w)
}

// Run the script registered on the called URL.
// Returns 200 if the check finished successfully.
// Returns 500 if the check failed or there is a runtime error.
func (app *application) runBashScript(w http.ResponseWriter, r *http.Request) {

	// Construct script path for execution
	script := app.scriptBase + r.URL.String() + ".sh"
	app.infoLog.Printf("Execute shell script: %s", script)

	// Execute bash script
	cmd := exec.Command("/bin/sh", script)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		// Check failed with defined message
		if out.String() != "" {
			app.infoLog.Printf("Script %s failed with output: %v", script, out.String())
			app.serverError(w, errors.New(out.String()))
			return
		}

		// Execution failed
		app.infoLog.Printf("Script %s finished with error: %v", script, stderr.String())
		app.serverError(w, errors.New(stderr.String()))
		return
	}

	app.infoLog.Println(out.String())
	app.serverTrue(w)
}
