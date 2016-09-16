package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"yotemplate"
)

// testPostPage is a struct for data on the post page
type testPostPage struct {
	Post          *Post
	OGURL         string
	OGType        string
	OGTitle       template.HTML
	OGDescription template.HTML
	OGImage       string
	UserName      string
}

// TestPostHandler is a handler for post page processing
type TestPostHandler struct {
	template *yotemplate.Template
	db       *sql.DB
}

// InitTestPostHandler initialize TestPostHandler struct
func InitTestPostHandler(db *sql.DB, templatesPath string) *TestPostHandler {

	pathes := make([]string, 2)
	pathes[0] = filepath.Join(templatesPath, "layout.html")
	pathes[1] = filepath.Join(templatesPath, "test_post.html")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Test post page template is initialized.")

	return &TestPostHandler{template: templ, db: db}
}

// TestPostHandler - handler for post page
func (h TestPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// validateURL
	res := TestPostURLPattern.FindStringSubmatch(r.URL.Path)
	if res == nil {
		http.NotFound(w, r)
		return
	}

	pp := &testPostPage{}

	// get cookie and userName
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		log.Printf("Error during cookie read on post page: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sessionID := cookie.Value
	user, err := DBGetUserBySessionID(h.db, sessionID)
	if err == nil {
		pp.UserName = user.Name
	}

	postID, err := strconv.Atoi(res[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pp.Post, err = DBGetPostByID(h.db, postID)
	if err != nil {
		log.Printf("Error during db query for post page: %v\n", err)
		http.NotFound(w, r)
		return
	}

	pp.OGDescription = pp.Post.Description
	pp.OGImage = pp.Post.ImageURL
	pp.OGTitle = pp.Post.Title
	pp.OGType = "article"
	pp.OGURL = "http://" + r.Host + r.URL.String()

	w.WriteHeader(http.StatusOK)

	h.template.Execute(w, pp)
}
