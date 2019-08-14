package main

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

func (app *application) startChecks() {

	for name, check := range app.checkList {
		app.infoLog.Printf("key: %s, value: %s", name, check)

		// Only run the check if active
		if check.active {
			go func(name string, check Check) {
				for {
					app.infoLog.Printf("Running check %s", name)

					result := app.runCheck(check)
					app.infoLog.Printf(result)

					value, labels := app.convertResult(result)
					app.infoLog.Printf("Result from check %s -> value: %s, labels: %v", name, value, labels)

					// TODO: Change metrics value and add labels

					time.Sleep(time.Duration(check.interval) * time.Second)
				}

			}(name, check)
		} else {
			app.infoLog.Printf("Check %s not active", check.name)
		}
	}
}

// Run the check and return the result.
func (app *application) runCheck(check Check) string {

	app.infoLog.Printf("Execute shell script: %s", check.file)

	// Execute bash script
	cmd := exec.Command("/bin/sh", check.file)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		// Check failed with defined message
		if out.String() != "" {
			app.infoLog.Printf("Script %s failed with output: %v", check.file, out.String())
			return out.String()
		}

		// Execution failed
		app.infoLog.Printf("Script %s finished with error: %v", check.file, stderr.String())
		return stderr.String()
	}

	// Check successfull
	app.infoLog.Printf("Script %s finished with success: %v", check.file, out.String())
	return out.String()
}

// Converts the return value from the script check.
// Format: value(int)|label1:value1,label2:value2
func (app *application) convertResult(result string) (string, map[string]string) {
	splitResult := strings.Split(result, "|")

	// Result of the check
	value := splitResult[0]

	// Labels of the check
	labels := make(map[string]string)
	splitLabels := strings.Split(splitResult[1], ",")
	for _, label := range splitLabels {
		splitLabel := strings.Split(label, ":")
		labels[splitLabel[0]] = splitLabel[1]
	}

	return value, labels
}
