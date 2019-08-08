package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	scriptBase string
}

func main() {

	// Parse command line paramters
	flagScriptBase := flag.String("scriptBase", "scripts", "Base path for the check scripts")
	flag.Parse()

	// Setup error handlers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog:   errorLog,
		infoLog:    infoLog,
		scriptBase: *flagScriptBase,
	}

	srv := &http.Server{
		Addr:     ":4444",
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Start the server
	infoLog.Printf("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
