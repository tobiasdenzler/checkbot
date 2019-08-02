package main

import (
	"net/http"
)

// Returns a false check result
func (app *application) serverError(w http.ResponseWriter, err error) {
	app.errorLog.Output(2, err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Returns a true check result
func (app *application) serverOk(w http.ResponseWriter) {
	w.Write([]byte("true"))
}
