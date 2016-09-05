package main

import (
	"log"
	"net/http"
	"path/filepath"
	"yotemplate"
)

// IndexPage is a struct for data on the index page
type IndexPage struct {
}

// IndexPageHandler is a handler for page processing
type IndexPageHandler struct {
	Template *yotemplate.YoTemplate
}

// InitIndexPageHandler initialize IndexPageHandler struct
func InitIndexPageHandler(templatesPath string) *IndexPageHandler {
	indexTemplatePath := filepath.Join(templatesPath, "index.html")
	indexTemplate, err := yotemplate.InitYoTemplate(indexTemplatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Index page template is initialized.")

	return &IndexPageHandler{Template: indexTemplate}
}

// IndexPageHandle - handler for index page
func (iph IndexPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

	}
	ip := &IndexPage{}

	w.WriteHeader(http.StatusOK)

	iph.Template.Execute(w, ip)
}
