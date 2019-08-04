package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	openshiftHost string
}

func main() {

	// Parse command line paramters
	addr := flag.String("addr", ":4444", "HTTP network address")
	openshiftHost := flag.String("openshiftHost", "127.0.0.1", "Hostname or IP of OpenShift API")
	flag.Parse()

	// Setup error handlers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		openshiftHost: *openshiftHost,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Start the server
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
