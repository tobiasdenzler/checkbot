package main

import (
	"net/http"
)

// Returns a false check result
func (app *application) serverError(w http.ResponseWriter, err error) {
	app.errorLog.Output(2, err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Returns a 200 and true check result
func (app *application) serverTrue(w http.ResponseWriter) {
	w.Write([]byte("true"))
}

// Returns a 404 and false check result
func (app *application) serverFalse(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("false"))
}
