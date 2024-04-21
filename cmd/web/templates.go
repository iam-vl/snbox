package main

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/iam-vl/snbox/internal/models"
)

// define templatedata as holding structure
// for all dynamic data to pass to HTML templates
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	// Get a slice of all filepaths that match "./ui/html/pages/*.tmpl"
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// Extract filename
		name := filepath.Base(page)
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}
		// Parse files into a template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func NTCache() (map[string]*template.Template, error) {
	fmt.Println("Starting ntcache... ")
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// Add base templ
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = template.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		cache[name] = ts

	}
	return cache, nil
}
