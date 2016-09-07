package main

import (
	"log"
	"path/filepath"
	"yotemplate"
)

// ErrorTemplate is a page for error mesages
var ErrorTemplate *yotemplate.YoTemplate

// InitErrorTemplate initialize error page template
func InitErrorTemplate(templatesPath string) *yotemplate.YoTemplate {
	templatePath := filepath.Join(templatesPath, "error.html")
	template, err := yotemplate.InitYoTemplate(templatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Error template is initialized.")
	return template
}
