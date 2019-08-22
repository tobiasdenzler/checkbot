package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Check consists of a script and a metric to scrape.
type Check struct {
	name        string
	file        string
	interval    int
	active      bool
	metricType  string
	help        string
	metric      interface{}
	stoppedchan chan struct{}
	nextrun     int64
}

// Define the metadata that can be used in the scripts
const metaActive = "ACTIVE"
const metaType = "TYPE"
const metaHelp = "HELP"
const metaInterval = "INTERVAL"

// Read all the available scripts and create a list of checks.
func (app *application) buildMetrics() {

	// Walk through all scripts and register the files with a handler
	err := filepath.Walk(app.scriptBase, func(path string, info os.FileInfo, err error) error {

		// Check if we have a file
		if info != nil && !info.IsDir() {

			// Openshift is using linking of files with ..
			if !strings.Contains(path, "..") {

				// Retrieve the status as bool
				active, _ := strconv.ParseBool(extractMetadataFromFile(metaActive, path))

				// Retrieve the interval as integer
				interval, _ := strconv.Atoi(extractMetadataFromFile(metaInterval, path))

				// Create a new check
				check := new(Check)
				check = &Check{
					name:        app.metricsPrefix + "_" + strings.Split(info.Name(), ".")[0], // Remove file ending
					file:        path,
					interval:    interval,
					active:      active,
					metricType:  extractMetadataFromFile(metaType, path),
					help:        extractMetadataFromFile(metaHelp, path),
					stoppedchan: make(chan struct{}),
					nextrun:     time.Now().Unix(),
				}

				// Add the check to the list
				app.checkList[check.name] = *check
				log.Infof("Add new check: %s", check.name)
				log.Debugf("Check details: %s", check.String())
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("Failed to read the scripts: %v", err)
	}
}

// Extract metadata information from a script.
// Metadata can be added using e.g. # TYPE
func extractMetadataFromFile(metadata string, file string) string {
	line, err := findLineInFile(file, "# "+metadata)
	if err == nil {
		return strings.TrimSpace(strings.Split(line, "# "+metadata)[1])
	}
	log.Errorf("Failed to retrieve %s from file %s", metadata, file)
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
		"[%s : %s : %d : %v : %s : %s]",
		c.name,
		c.file,
		c.interval,
		c.active,
		c.metricType,
		c.help)
}
