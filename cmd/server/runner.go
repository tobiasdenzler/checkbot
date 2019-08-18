package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

// Starts a go routine for each check in the list.
func (app *application) startChecks() {

	// Walk throught the check list
	for name, check := range app.checkList {
		log.Debugf("key: %s, value: %s", name, check)

		// Only run the check if active
		if check.active {

			// Start new go routine
			go func(name string, check Check) {
				for {
					log.Debugf("Running check %s", name)

					// Run the script
					result, err := runCheck(check)

					if err == nil {

						// Split the result from the check script, can be multiple lines
						resultLine := strings.Split(result, "\n")
						for _, line := range resultLine {
							if line != "" {
								value, labels := convertResult(line)

								// TODO: Support other type of metrics
								switch check.metricType {
								case "Gauge":
									if check.metric == nil {
										check.metric = prometheus.NewGaugeVec(
											prometheus.GaugeOpts{
												Name: check.name,
												Help: check.help,
											},
											convertMapKeysToSlice(labels),
										)
										prometheus.MustRegister(check.metric.(*prometheus.GaugeVec))
									}
									check.metric.(*prometheus.GaugeVec).With(labels).Set(value)
								default:
									check.metric = nil
								}

								log.Debugf("Result from check %s -> value: %f, labels: %v", name, value, labels)
							}
						}
					} else {
						log.Warnf("Check %s failed with error: %s", name, err)
					}

					// Wait for the defined interval
					time.Sleep(time.Duration(check.interval) * time.Second)
				}

			}(name, check)
		} else {
			log.Infof("Check %s not active", check.name)
		}
	}
}

// Run the check and return the result.
func runCheck(check Check) (string, error) {

	log.Debugf("Execute shell script: %s", check.file)

	// Execute bash script
	cmd := exec.Command("/bin/sh", check.file)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		// Check failed with defined message
		if out.String() != "" {
			log.Infof("Script %s failed with output: %v", check.file, out.String())
			return "", errors.New("Script failed with error: " + out.String())
		}

		// Execution failed
		log.Infof("Script %s finished with error: %v", check.file, stderr.String())
		return "", errors.New("Script failed with error: " + stderr.String())
	}

	// Check has error
	if out.String() == "" {
		log.Infof("Script %s finished with error: %v", check.file, stderr.String())
		return "", errors.New("Script failed with error: " + stderr.String())
	}

	// Check run successfull
	log.Debugf("Script %s finished with success: %v", check.file, out.String())
	return out.String(), nil
}

// Converts the return value from the script check.
// Format: value|label1:value1,label2:value2
func convertResult(result string) (float64, map[string]string) {
	var metricValue float64
	var labels = make(map[string]string)

	if strings.Contains(result, "|") {
		splitResult := strings.Split(result, "|")

		// Result of the check
		value := splitResult[0]

		// Labels of the check
		splitLabels := strings.Split(splitResult[1], ",")
		for _, label := range splitLabels {
			splitLabel := strings.SplitN(label, ":", 2)
			labels[splitLabel[0]] = splitLabel[1]
		}
		metricValue, _ = strconv.ParseFloat(value, 64)
	} else {
		metricValue, _ = strconv.ParseFloat(result, 64)
	}
	return metricValue, labels
}

// Convert the keys from a map to a slice.
func convertMapKeysToSlice(value map[string]string) []string {
	keys := make([]string, len(value))

	i := 0
	for k := range value {
		keys[i] = k
		i++
	}

	return keys
}
