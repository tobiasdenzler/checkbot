package main

import (
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct {
	Checklist     map[string]*Check
	Sandbox       Sandbox
	Configuration Configuration
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	// Initialize new map to act as cache
	cache := map[string]*template.Template{}

	// Get all 'page' templates for the application
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract file name from the full file path
		name := filepath.Base(page)

		// Parse page template file in to template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add any 'layout' templates to thetemplate set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add any 'partial' templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add template set to the cache, using name of the page as key
		cache[name] = ts
	}

	return cache, nil
}

// Converts a unix timestamp to human readable datetime
func humanDate(t int64) string {
	current := time.Unix(t, int64(0))
	return current.Format("2006-01-02 15:04:05")
}

// Initialize a template.FuncMap object and store it in a global variable.
var functions = template.FuncMap{
	"humanDate": humanDate,
}
