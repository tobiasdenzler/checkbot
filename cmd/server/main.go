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
	logLevel      string
	checkList     map[string]Check
	templateCache map[string]*template.Template
}

func init() {
	// set default log config
	log.SetLevel(log.WarnLevel)
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {

	// Parse command line paramters
	flagScriptBase := flag.String("scriptBase", "scripts", "Base path for the check scripts")
	flagMetricsPrefix := flag.String("metricsPrefix", "checkbot", "Prefix for all metrics")
	flagLogLevel := flag.String("logLevel", "info", "Log level for application (error|warn|info|debug|trace")
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
		logLevel:      *flagLogLevel,
		checkList:     checkList,
		templateCache: templateCache,
	}

	// parse custom loglevel
	level, err := log.ParseLevel(*flagLogLevel)

	// set loglevel based on config
	log.SetLevel(level)
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Build metrics and fill checklist
	app.buildMetrics()

	// Start running the checks
	app.startChecks()

	// Start the server
	log.Infof("Starting server on :4444")
	err = http.ListenAndServeTLS(":4444", "./certs/server.crt", "./certs/server.key", app.routes())
	log.Fatal(err)
}
