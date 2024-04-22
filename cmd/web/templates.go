package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/iam-vl/snbox/internal/models"
)

// define templatedata as holding structure
// for all dynamic data to pass to HTML templates
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func HumanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var Functions = template.FuncMap{
	"humanDate": HumanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// Extract filename
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(Functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}
		// files := []string{
		// 	"./ui/html/base.tmpl",
		// 	"./ui/html/partials/nav.tmpl",
		// 	page,
		// }
		// ts, err := template.ParseFiles(files...)
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
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
