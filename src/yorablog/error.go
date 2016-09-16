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
	templatePath := filepath.Join(templatesPath, "error.html")
	template, err := yotemplate.InitTemplate(templatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Error template is initialized.")
	return template
}
