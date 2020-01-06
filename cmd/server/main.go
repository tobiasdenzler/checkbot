package main

import (
	"flag"
	"html/template"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Global data
type application struct {
	scriptBase    string
	metricsPrefix string
	logLevel      string
	managementPwd string
	enableSandbox bool
	checkList     map[string]*Check
	lastrunMetric *prometheus.GaugeVec
	templateCache map[string]*template.Template
	config        Configuration
}

func init() {
	// set default log config
	log.SetLevel(log.WarnLevel)
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

// Version is provided by ldflags
var Version = "unspecified"

// Build is provided by ldflags
var Build = "unspecified"

func main() {
	// Parse command line paramters
	flagScriptBase := flag.String("scriptBase", "scripts", "Base path for the check scripts")
	flagMetricsPrefix := flag.String("metricsPrefix", "checkbot", "Prefix for all metrics")
	flagLogLevel := flag.String("logLevel", "info", "Log level for application (error|warn|info|debug|trace")
	flagManagementPwd := flag.String("managementPwd", "admin", "Password for managing endpoints")
	flagEnableSandbox := flag.Bool("enableSandbox", false, "Enable debugging sandbox")
	flag.Parse()

	// Create map for all checks
	checkList := map[string]*Check{}

	// Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize config values
	config := &Configuration{
		Version: Version,
		Build:   Build,
		Sandbox: *flagEnableSandbox,
	}

	// Global application variables
	app := &application{
		scriptBase:    *flagScriptBase,
		metricsPrefix: *flagMetricsPrefix,
		logLevel:      *flagLogLevel,
		managementPwd: *flagManagementPwd,
		checkList:     checkList,
		lastrunMetric: nil,
		templateCache: templateCache,
		config:        *config,
	}

	// parse custom loglevel
	level, err := log.ParseLevel(*flagLogLevel)

	// set loglevel based on config
	log.SetLevel(level)
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Show build information
	log.Infof("Version: %s, Build: %s", Version, Build)

	// Build metrics and fill checklist
	app.buildMetrics()

	// Start running the checks
	app.startChecks()

	// Start the server
	log.Infof("Starting server on :4444")
	err = http.ListenAndServeTLS(":4444", "./certs/tls.crt", "./certs/tls.key", app.routes())
	log.Fatal(err)
}
