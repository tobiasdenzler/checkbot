package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Global data
type application struct {
	scriptBase    string
	metricsPrefix string
	checkList     map[string]Check
}

func init() {
	// set loglevel based on config
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {

	// Parse command line paramters
	flagScriptBase := flag.String("scriptBase", "scripts", "Base path for the check scripts")
	flagMetricsPrefix := flag.String("metricsPrefix", "checkbot", "Prefix for all metrics")
	flag.Parse()

	// Create map for all checks
	checkList := map[string]Check{}

	app := &application{
		scriptBase:    *flagScriptBase,
		metricsPrefix: *flagMetricsPrefix,
		checkList:     checkList,
	}

	srv := &http.Server{
		Addr:    ":4444",
		Handler: app.routes(),
	}

	// Build metrics and fill checklist
	app.buildMetrics()

	// Start running the checks
	app.startChecks()

	// Start the server
	log.Infof("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
