package main

import (
	"html/template"
	"path/filepath"
)

type templateData struct {
	Checklist map[string]Check
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
		ts, err := template.ParseFiles(page)
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
