package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"yotemplate"
)

// EditPage is a struct for data on the edit page
type EditPage struct {
	Post     *Post
	UserName string
}

// EditPageHandler is a handler for edit page processing
type EditPageHandler struct {
	Template *yotemplate.Template
}

// InitEditPageHandler initialize EditPageHandler struct
func InitEditPageHandler(templatesPath string) *EditPageHandler {
	editTemplatePath := filepath.Join(templatesPath, "edit.html")
	editTemplate, err := yotemplate.InitTemplate(editTemplatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Edit page template is initialized.")

	return &EditPageHandler{Template: editTemplate}
}

// EditPageHandle - handler for index page
func (eph EditPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	res := EditURLPattern.FindStringSubmatch(r.URL.Path)
	if res == nil {
		http.NotFound(w, r)
		return
	}

	postID, err := strconv.Atoi(res[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		ep := &EditPage{}

		cookie, err := r.Cookie("SessionID")
		if err != nil {
			log.Printf("Error during cookie read on index page: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sessionID := cookie.Value
		user, err := DBGetUserBySessionID(DBConnection, sessionID)
		if err == nil {
			ep.UserName = user.Name
		}

		if !user.EditPostPermit {
			w.WriteHeader(http.StatusForbidden)
			ErrorTemplate.Execute(w, "Недостаточно прав для редактирования статьи")
			return
		}

		ep.Post, err = DBGetPostByID(DBConnection, postID)
		if err != nil {
			log.Printf("Error during db query for edit page: %v\n", err)
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)

		eph.Template.Execute(w, ep)
	} else if r.Method == "POST" {
		err = r.ParseForm()
		if err != nil {
			log.Printf("Error during edit form parse: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		post := &Post{ID: postID}

		post.Title = template.HTML(r.FormValue("title"))
		post.Description = template.HTML(r.FormValue("description"))
		post.ImageURL = r.FormValue("imageurl")
		post.Annotation = template.HTML(r.FormValue("annotation"))
		post.Text = template.HTML(r.FormValue("posttext"))

		err = DBUpdatePost(post)
		if err != nil {
			log.Printf("Error during update post in DB: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		redirectURL := fmt.Sprintf("/post/%v", post.ID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
