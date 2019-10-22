package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Error(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	// Retrieve appropriate template set from the cache
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// Write to buffer first
	buf := new(bytes.Buffer)

	// Execute the template set, passing in any dynamic data
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

// MapToString will convert a map to a string.
func MapToString(m map[string]string) string {
	tmp := ""
	for key, value := range m {
		tmp += key + ":" + value + ","
	}
	return strings.TrimSuffix(tmp, ",")
}
