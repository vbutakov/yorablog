package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"yotemplate"
)

// CreatePage is a struct for data on the create page
type CreatePage struct {
	Post *Post
}

// CreatePageHandler is a handler for edit create processing
type CreatePageHandler struct {
	Template *yotemplate.YoTemplate
}

// InitCreatePageHandler initialize CreatePageHandler struct
func InitCreatePageHandler(templatesPath string) *CreatePageHandler {
	createTemplatePath := filepath.Join(templatesPath, "create.html")
	createTemplate, err := yotemplate.InitYoTemplate(createTemplatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Create page template is initialized.")

	return &CreatePageHandler{Template: createTemplate}
}

// CreatePageHandle - handler for create page
func (cph CreatePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		cp := &CreatePage{Post: &Post{}}

		w.WriteHeader(http.StatusOK)

		cph.Template.Execute(w, cp)
	} else if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Printf("Error during create page form parse: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		post := &Post{}

		post.Title = template.HTML(r.FormValue("title"))
		post.Description = template.HTML(r.FormValue("description"))
		post.ImageURL = r.FormValue("imageurl")
		post.Annotation = template.HTML(r.FormValue("annotation"))
		post.Text = template.HTML(r.FormValue("posttext"))

		var postID int
		postID, err = DBInsertPost(post)
		if err != nil {
			log.Printf("Error during insert post in DB: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		redirectURL := fmt.Sprintf("/post/%v", postID)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}
