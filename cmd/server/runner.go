package main

import (
	"bytes"
	"errors"
	"os/exec"
	"reflect"
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

	app.registerLastrunMetric()
	app.registerLastresultMetric()

	log.Debug("Starting all checks now..")

	// Recreate the chan in case it was closed before
	stopchan = make(chan struct{})

	// Walk throught the check list
	for _, check := range app.checkList {
		// Only run the check if active
		if check.Active {
			go app.runCheck(check, stopchan)
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

	// Reset the status metrics
	prometheus.Unregister(app.lastresultMetric)
	log.Debug("Unregistered lastresult metric")
	prometheus.Unregister(app.lastrunMetric)
	log.Debug("Unregistered lastrun metric")

	log.Debug("All checks are stopped.")
}

// Run the check and save the result to the list.
func (app *application) runCheck(check *Check, stopchan chan struct{}) {

	// Close the stoppedchan when this func exits
	defer close(check.stoppedchan)

	// Teardown
	defer func() {
		unregisterMetricsForCheck(check)
	}()

	for {
		select {
		default:

			// Check if we can run the check
			if time.Now().Unix() > check.Nextrun {

				log.Debugf("Running check %s", check.Name)

				// Store result of previous run
				check.resultLast = check.resultCurrent
				check.resultCurrent = []map[string]string{}

				// Run the script
				result, err := runBashScript(*check)

				check.Success = 0
				if err == nil {
					check.Success = 1

					// Split the result from the check script, can be multiple lines
					resultLine := strings.Split(result, "\n")
					for _, line := range resultLine {
						if line != "" {
							// Extract values from the result and register the metric
							value, labels := convertResult(line)
							registerMetricsForCheck(check, value, labels)
						}
					}

				} else {
					log.Warnf("Check %s failed with error: %s", check.Name, err)
				}

				// Cleanup stale metrics data
				cleanupUnusedDimensions(check)

				// Set time for next run
				check.Nextrun = check.Nextrun + int64(check.Interval) + check.Offset
				log.Debugf("Finished check %s and schedule next run for %s", check.Name, time.Unix(check.Nextrun, 0))

				// Update lastrun metric
				lastStatusLabels := make(map[string]string)
				lastStatusLabels["name"] = check.Name
				lastStatusLabels["interval"] = strconv.Itoa(check.Interval)
				lastStatusLabels["offset"] = strconv.FormatInt(check.Offset, 10)
				lastStatusLabels["type"] = check.MetricType

				app.lastrunMetric.With(lastStatusLabels).Set(float64(time.Now().Unix()))
				app.lastresultMetric.With(lastStatusLabels).Set(float64(check.Success))

				log.Debugf("lastresult is %v", check.Success)
				log.Debugf("Adding lastStatusLabels for %s with values %v", check.Name, lastStatusLabels)
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

// Register all metrics from Prometheus for a given check.
func registerMetricsForCheck(check *Check, value float64, labels map[string]string) {

	defer func() {
		if r := recover(); r != nil {
			log.Warnf("Not able to register metric for check %s. Maybe already registered?", check.Name)
		}
	}()

	// Store the result labels
	check.resultCurrent = append(check.resultCurrent, labels)

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

			// This can be panicking and will be recovered
			prometheus.MustRegister(check.metric.(*prometheus.GaugeVec))
		}
		check.metric.(*prometheus.GaugeVec).With(labels).Set(value)
	case "Counter":
		if check.metric == nil {
			check.metric = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: check.Name,
					Help: check.Help,
				},
				convertMapKeysToSlice(labels),
			)

			// This can be panicking and will be recovered
			prometheus.MustRegister(check.metric.(*prometheus.CounterVec))
		}
		check.metric.(*prometheus.CounterVec).With(labels).Add(value)
	case "Histogram":
		log.Warn("Metric type Counter not implemented yet!")
	case "Summary":
		log.Warn("Metric type Counter not implemented yet!")
	default:
		log.Warnf("Not able to register unknown metric type %s", check.MetricType)
		check.metric = nil
	}

	log.Tracef("Result from check %s -> value: %f, labels: %v", check.Name, value, MapToString(labels))
}

// Cleanup metric vectors we do not need anymore.
func cleanupUnusedDimensions(check *Check) {

	log.Tracef("Check %s cleaning up -> size of resultLast : %d, size of resultCurrent: %d", check.Name, len(check.resultLast), len(check.resultCurrent))

	// Loop through labels from last run and check if they are still valid for
	// the current run, otherwise remove them.
	var remove bool
	for _, labelsLast := range check.resultLast {
		remove = true

		if len(check.resultCurrent) > 0 {
			for _, labelCurrent := range check.resultCurrent {
				if reflect.DeepEqual(labelsLast, labelCurrent) {
					remove = false
				}
			}
		}

		// Remove the stale metric
		if remove {
			log.Debugf("Check %s remove stale metric vector with labels %s", check.Name, MapToString(labelsLast))

			switch check.MetricType {
			case "Gauge":
				if !(check.metric.(*prometheus.GaugeVec).Delete(labelsLast)) {
					log.Warnf("Failed to delete stale metric vector with label %s from check %s", MapToString(labelsLast), check.Name)
				}
			case "Counter":
				if !(check.metric.(*prometheus.CounterVec).Delete(labelsLast)) {
					log.Warnf("Failed to delete stale metric vector with label %s from check %s", MapToString(labelsLast), check.Name)
				}
			default:
				log.Warnf("Not able to remove unknown metric type %s", check.MetricType)
			}
		}
	}
}

// Unregister all metrics from Prometheus for a given check.
func unregisterMetricsForCheck(check *Check) {
	if check.metric != nil {
		switch check.MetricType {
		case "Gauge":
			prometheus.Unregister(check.metric.(*prometheus.GaugeVec))
		case "Counter":
			prometheus.Unregister(check.metric.(*prometheus.CounterVec))
		case "Histogram":
			log.Warn("Metric type Counter not implemented yet!")
		case "Summary":
			log.Warn("Metric type Counter not implemented yet!")
		default:
			log.Warnf("Not able to unregister unknown metric type %s", check.MetricType)
		}
		check.metric = nil

		log.Debugf("Unregistered metrics for check %s", check.Name)
	}
}

// Run the check and return the result.
func runBashScript(check Check) (string, error) {

	log.Debugf("Execute shell script: %s", check.File)

	// Execute bash script
	cmd := exec.Command(check.File)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	scriptResult := out.String()
	scriptError := stderr.String()

	if err != nil {
		// Check failed with defined message
		if scriptResult != "" {
			log.Infof("Script %s failed with output: %v", check.File, scriptResult)
			return "", errors.New("Script failed with error: " + scriptResult)
		}

		// Check has error
		if scriptError != "" {
			log.Infof("Script %s failed with error: %v", check.File, scriptError)
			return "", errors.New("Script failed with error: " + scriptError)
		}

		// Execution failed
		log.Infof("Script %s finished with execution error: %v", check.File, err)
		return "", errors.New("Script failed with error: " + err.Error())
	}

	// Check run successfull
	return scriptResult, nil
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
		log.Tracef("--> split labels: %v", splitLabels)
		for _, label := range splitLabels {
			if(!strings.Contains(label, "=")) {
				log.Warnf("Skipping label %s because wrong format detected in result: %s", label, result)
			} else {
				splitLabel := strings.SplitN(label, "=", 2)
				labels[splitLabel[0]] = splitLabel[1]
			}
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

// Setup the lastrun metric for information about the last execution time of a checks
func (app *application) registerLastrunMetric() {
	app.lastrunMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checkbot_lastrun_info",
			Help: "Provides information about the last run of a script.",
		},
		[]string{"name", "interval", "offset", "type"},
	)

	// Metric could already be registered, but this is not a problem
	prometheus.Register(app.lastrunMetric)
	log.Debug("Registering metric lastrun")
}

// Setup the lastresult metric for information about the last result of a checks
func (app *application) registerLastresultMetric() {
	app.lastresultMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checkbot_lastresult_info",
			Help: "Provides information about the last result of a script.",
		},
		[]string{"name", "interval", "offset", "type"},
	)

	// Metric could already be registered, but this is not a problem
	prometheus.Register(app.lastresultMetric)
	log.Debug("Registering metric lastresult")
}
