package main

import (
	"embed"
	_ "embed"
	"html/template"
	"path/filepath"
)

//go:embed ui
var content embed.FS

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages := []string{
		"ui/html/pages/home.html",
		"ui/html/pages/board.html",
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"ui/html/base.html",
			page,
		}

		ts, err := template.ParseFS(content, files...)
		// ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
