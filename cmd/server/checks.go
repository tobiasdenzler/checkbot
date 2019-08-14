package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Check consists of a script and a metrics to scrape.
type Check struct {
	name       string
	file       string
	interval   int
	active     bool
	metricType string
	help       string
	metric     interface{}
}

// Define the metadata that can be used in the scripts
const metaActive = "ACTIVE"
const metaType = "TYPE"
const metaHelp = "HELP"
const metaInterval = "INTERVAL"

// Define the routes to serve
func (app *application) buildMetrics() {

	// Walk through all scripts and register the files with a handler
	err := filepath.Walk(app.scriptBase, func(path string, info os.FileInfo, err error) error {

		// Check if we have a file
		if info != nil && !info.IsDir() {

			// Openshift is using linking of files with ..
			if !strings.Contains(path, "..") {

				// Retrieve the status as bool
				active, _ := strconv.ParseBool(app.extractMetadataFromFile(metaActive, path))
				
				// Retrieve the interval as integer
				interval, _ := strconv.Atoi(app.extractMetadataFromFile(metaInterval, path))

				// Create a new check
				check := new(Check)
				check = &Check{
					name:       app.metricsPrefix + "_" + strings.Split(info.Name(), ".")[0], // Remove file ending
					file:       path,
					interval:   interval,
					active:     active,
					metricType: app.extractMetadataFromFile(metaType, path),
					help:       app.extractMetadataFromFile(metaHelp, path),
				}

				// TODO: Support other type of metrics
				metric := promauto.NewGauge(prometheus.GaugeOpts{
					Name: check.name,
					Help: check.help,
					//ConstLabels: map[string]string{"project": "test"},
				})
				check.metric = metric

				app.checkList[check.name] = *check
				app.infoLog.Printf("Add new check %s", check.String())
			}
		}
		return nil
	})
	if err != nil {
		app.errorLog.Printf("Failed to read the scripts: %v", err)
	}
}

// Extract metadata information from a script.
// Metadata can be added using e.g. # TYPE
func (app *application) extractMetadataFromFile(metadata string, file string) string {
	line, err := findLineInFile(file, "# "+metadata)
	if err == nil {
		return strings.TrimSpace(strings.Split(line, "# "+metadata)[1])
	}
	app.errorLog.Printf("Failed to retrieve %s from file %s", metadata, file)
	return ""
}

// Search for a string in a file and return the corresponding line.
func findLineInFile(path string, searchFor string) (string, error) {

	// Open file for reading
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchFor) {
			return scanner.Text(), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("Failed to find " + searchFor + " in file " + path + ".")
}

// String returns the Check as string.
func (c Check) String() string {
	return fmt.Sprintf(
		"[%s : %s : %d : %s : %s]",
		c.name,
		c.file,
		c.interval,
		c.metricType,
		c.help)
}
