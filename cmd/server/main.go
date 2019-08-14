package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Global data
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	scriptBase    string
	metricsPrefix string
	checkList     map[string]Check
}

func main() {

	// Parse command line paramters
	flagScriptBase := flag.String("scriptBase", "scripts", "Base path for the check scripts")
	flagMetricsPrefix := flag.String("metricsPrefix", "checkbot", "Prefix for all metrics")
	flag.Parse()

	// Setup error handlers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create map for all checks
	checkList := map[string]Check{}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		scriptBase:    *flagScriptBase,
		metricsPrefix: *flagMetricsPrefix,
		checkList:     checkList,
	}

	srv := &http.Server{
		Addr:     ":4444",
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Build metrics and fill checklist
	app.buildMetrics()

	// Start running the checks
	app.startChecks()

	// Start the server
	infoLog.Printf("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
