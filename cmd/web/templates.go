package main

import "github.com/iam-vl/snbox/internal/models"

// define templatedata as holding structure
// for all dynamic data to pass to HTML templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
