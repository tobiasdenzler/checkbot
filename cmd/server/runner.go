package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Starts a go routine for each check in the list.
func (app *application) startChecks() {

	// Walk throught the check list
	for name, check := range app.checkList {
		app.infoLog.Printf("key: %s, value: %s", name, check)

		// Only run the check if active
		if check.active {

			// Start new go routine
			go func(name string, check Check) {
				for {
					app.infoLog.Printf("Running check %s", name)

					// Run the script
					result := app.runCheck(check)

					// Split the result from the check script, can be multiple lines
					resultLine := strings.Split(result, "\n")
					for _, line := range resultLine {
						if line != "" {
							value, labels := convertResult(line)

							// TODO: Support other type of metrics
							switch check.metricType {
							case "Gauge":
								if check.metric == nil {
									check.metric = promauto.NewGauge(prometheus.GaugeOpts{
										Name:        check.name,
										Help:        check.help,
										ConstLabels: labels,
									})
								}
								check.metric.(prometheus.Gauge).Set(value)
							default:
								check.metric = nil
							}

							app.infoLog.Printf("Result from check %s -> value: %f, labels: %v", name, value, labels)
						}
					}

					// Wait for the defined interval
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
// Format: value|label1:value1,label2:value2
func convertResult(result string) (float64, map[string]string) {
	var metricValue float64
	var labels map[string]string

	if strings.Contains(result, "|") {
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
		metricValue, _ = strconv.ParseFloat(value, 64)
	} else {
		metricValue, _ = strconv.ParseFloat(result, 64)
	}
	return metricValue, labels
}
