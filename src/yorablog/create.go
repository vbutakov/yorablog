package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"yoradb"
	"yotemplate"
)

// CreatePage is a struct for data on the create page
type CreatePage struct {
	Post         *yoradb.Post
	UserName     string
	ErrorMessage string
}

// CreatePageHandler is a handler for edit create processing
type CreatePageHandler struct {
	template *yotemplate.Template
	db       yoradb.PostRepository
}

// InitCreatePageHandler initialize CreatePageHandler struct
func InitCreatePageHandler(db yoradb.PostRepository, templatesPath string) *CreatePageHandler {

	pathes := make([]string, 3)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "create.gohtml")
	pathes[2] = filepath.Join(templatesPath, "empty_og.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Create page template is initialized.")

	return &CreatePageHandler{template: templ, db: db}
}

// CreatePageHandle - handler for create page
func (h CreatePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cp := &CreatePage{Post: &yoradb.Post{}}

	user, ok := UserFromContext(r.Context())
	if ok {
		cp.UserName = user.Name
	}

	if !user.CreatePostPermit {
		w.WriteHeader(http.StatusForbidden)
		cp.ErrorMessage = "Недостаточно прав для создания статьи"
		ErrorTemplate.Execute(w, cp)
		return
	}

	if r.Method == http.MethodGet {

		w.WriteHeader(http.StatusOK)

		h.template.Execute(w, cp)

	} else if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			log.Printf("Error during create page form parse: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		post := &yoradb.Post{}

		post.Title = template.HTML(r.FormValue("title"))
		post.Description = template.HTML(r.FormValue("description"))
		post.ImageURL = r.FormValue("imageurl")
		post.Annotation = template.HTML(r.FormValue("annotation"))
		post.Text = template.HTML(r.FormValue("posttext"))

		var postID int64
		postID, err = h.db.CreatePost(post, user.ID)
		if err != nil {
			cp.Post = post
			cp.ErrorMessage = err.Error()

			log.Printf("Error during insert post in DB: %v\nUserID: %v\n", err, user.ID)

			w.WriteHeader(http.StatusOK)
			h.template.Execute(w, cp)

			return
		}

		redirectURL := fmt.Sprintf("/post/%v", postID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
