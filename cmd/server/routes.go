package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Define the routes to serve
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Internal health check
	mux.HandleFunc("/health", app.health)

	// Walk through all scripts and register the files with a handler
	err := filepath.Walk(app.scriptBase, func(path string, info os.FileInfo, err error) error {
		// Check if we have a file
		if !info.IsDir() {
			if !strings.Contains(path, "..") {
				// Construct the path of the script
				scriptPath := "/" + strings.Split(strings.Replace(path, app.scriptBase+"/", "", 1), ".")[0]
				// Register the path
				mux.HandleFunc(scriptPath, app.runBashScript)
				app.infoLog.Printf("registering path %s with script %s", scriptPath, path)
			}
		}
		return nil
	})
	if err != nil {
		app.errorLog.Panic(err)
	}

	return mux
}
