package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// A channel to tell it to stop
var stopchan chan struct{}

// Starts a go routine for each check in the list.
func (app *application) startChecks() {

	log.Debug("Starting all checks now..")

	// Recreate the chan in case it was closed before
	stopchan = make(chan struct{})

	// Walk throught the check list
	for _, check := range app.checkList {
		// Only run the check if active
		if check.Active {
			go runCheck(check, stopchan)
		} else {
			log.Infof("Check %s not active", check.Name)
		}
	}
}

// Stop all running go routines.
func (app *application) stopChecks() {

	log.Debug("Stopping all checks now..")
	close(stopchan)

	// Walk throught the check list
	for _, check := range app.checkList {
		if check.Active {
			<-check.stoppedchan
		}
	}
	log.Debug("All checks are stopped.")
}

// Run the check and save the result to the list.
func runCheck(check Check, stopchan chan struct{}) {

	// Close the stoppedchan when this func exits
	defer close(check.stoppedchan)

	// Teardown
	defer func() {

		// Unregister the metrics
		switch check.MetricType {
		case "Gauge":
			prometheus.Unregister(check.metric.(*prometheus.GaugeVec))
		}
	}()

	for {
		select {
		default:

			// Check if we can run the check
			if time.Now().Unix() > check.nextrun {

				log.Debugf("Running check %s", check.Name)

				// Run the script
				result, err := runBashScript(check)

				if err == nil {
					// Split the result from the check script, can be multiple lines
					resultLine := strings.Split(result, "\n")
					for _, line := range resultLine {
						if line != "" {
							value, labels := convertResult(line)

							// TODO: Support other type of metrics
							switch check.MetricType {
							case "Gauge":
								if check.metric == nil {
									check.metric = prometheus.NewGaugeVec(
										prometheus.GaugeOpts{
											Name: check.Name,
											Help: check.Help,
										},
										convertMapKeysToSlice(labels),
									)
									prometheus.MustRegister(check.metric.(*prometheus.GaugeVec))
								}
								check.metric.(*prometheus.GaugeVec).With(labels).Set(value)
							default:
								check.metric = nil
							}

							log.Tracef("Result from check %s -> value: %f, labels: %v", check.Name, value, labels)
						}
					}
				} else {
					log.Warnf("Check %s failed with error: %s", check.Name, err)
				}

				// Set time for next run
				check.nextrun += int64(check.Interval)
			}

		case <-stopchan:
			// Stop
			log.Debugf("Stopping check %s", check.Name)
			return

		case <-time.After(10 * time.Second):
			// Task didn't stop in time
			log.Debugf("Forced stopping check %s", check.Name)
			return
		}

		// Slow down
		time.Sleep(1 * time.Second)
	}
}

// Run the check and return the result.
func runBashScript(check Check) (string, error) {

	log.Debugf("Execute shell script: %s", check.File)

	// Execute bash script
	cmd := exec.Command("/bin/sh", check.File)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		// Check failed with defined message
		if out.String() != "" {
			log.Infof("Script %s failed with output: %v", check.File, out.String())
			return "", errors.New("Script failed with error: " + out.String())
		}

		// Execution failed
		log.Infof("Script %s finished with error: %v", check.File, stderr.String())
		return "", errors.New("Script failed with error: " + stderr.String())
	}

	// Check has error
	if out.String() == "" {
		log.Infof("Script %s finished with error: %v", check.File, stderr.String())
		return "", errors.New("Script failed with error: " + stderr.String())
	}

	// Check run successfull
	log.Tracef("Script %s finished with success: %v", check.File, out.String())
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
			splitLabel := strings.SplitN(label, "=", 2)
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
