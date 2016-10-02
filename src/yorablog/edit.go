package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"yoradb"
	"yotemplate"
)

// EditPage is a struct for data on the edit page
type EditPage struct {
	Post         *yoradb.Post
	UserName     string
	ErrorMessage string
}

// EditPageHandler is a handler for edit page processing
type EditPageHandler struct {
	template *yotemplate.Template
	db       yoradb.PostRepository
}

// InitEditPageHandler initialize EditPageHandler struct
func InitEditPageHandler(db yoradb.PostRepository, templatesPath string) *EditPageHandler {

	pathes := make([]string, 3)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "edit.gohtml")
	pathes[2] = filepath.Join(templatesPath, "empty_og.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Edit page template is initialized.")

	return &EditPageHandler{template: templ, db: db}
}

// EditPageHandle - handler for index page
func (h EditPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	res := EditURLPattern.FindStringSubmatch(r.URL.Path)
	if res == nil {
		http.NotFound(w, r)
		return
	}

	postID, err := strconv.ParseInt(res[1], 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ep := &EditPage{}

	user, ok := UserFromContext(r.Context())
	if ok {
		ep.UserName = user.Name
	}

	if !user.EditPostPermit {
		w.WriteHeader(http.StatusForbidden)
		ep.ErrorMessage = "Недостаточно прав для редактирования статьи"
		ErrorTemplate.Execute(w, ep)
		return
	}

	if r.Method == http.MethodGet {

		ep.Post, err = h.db.GetPostByID(postID)
		if err != nil {
			log.Printf("Error during db query for edit page: %v\n", err)
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)

		h.template.Execute(w, ep)
	} else if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			log.Printf("Error during edit form parse: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		post := &yoradb.Post{ID: postID}

		post.Title = template.HTML(r.FormValue("title"))
		post.Description = template.HTML(r.FormValue("description"))
		post.ImageURL = r.FormValue("imageurl")
		post.Annotation = template.HTML(r.FormValue("annotation"))
		post.Text = template.HTML(r.FormValue("posttext"))

		err = h.db.UpdatePost(post)
		if err != nil {
			ep.Post = post
			ep.ErrorMessage = err.Error()

			log.Printf("Error during update post in DB: %v\n", err)

			w.WriteHeader(http.StatusOK)
			h.template.Execute(w, ep)

			return
		}

		redirectURL := fmt.Sprintf("/post/%v", post.ID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
