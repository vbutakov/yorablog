package main

import (
	"log"
	"path/filepath"
	"yotemplate"
)

// ErrorTemplate is a page for error mesages
var ErrorTemplate *yotemplate.Template

// InitErrorTemplate initialize error page template
func InitErrorTemplate(templatesPath string) *yotemplate.Template {

	pathes := make([]string, 3)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "error.gohtml")
	pathes[2] = filepath.Join(templatesPath, "empty_og.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Error template is initialized.")
	return templ
}
