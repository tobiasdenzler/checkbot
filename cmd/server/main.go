package main

import (
	"flag"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Global data
type application struct {
	scriptBase    string
	metricsPrefix string
	checkList     map[string]Check
	templateCache map[string]*template.Template
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

	// Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		scriptBase:    *flagScriptBase,
		metricsPrefix: *flagMetricsPrefix,
		checkList:     checkList,
		templateCache: templateCache,
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
	err = srv.ListenAndServe()
	log.Fatal(err)
}
