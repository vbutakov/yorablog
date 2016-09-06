package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"yotemplate"
)

// IndexPage is a struct for data on the index page
type IndexPage struct {
	CurrentPage int
	Posts       []Post
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

	ok := IndexURLPattern.MatchString(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	var offset int
	var pageNum int
	var err error

	pageNumStr := IndexURLPattern.SubexpNames()[1]
	if pageNumStr == "" {
		offset = 0
	} else {
		pageNum, err = strconv.Atoi(pageNumStr)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if pageNum < 2 {
			http.NotFound(w, r)
			return
		}
		offset = (pageNum - 1) * 10 // 10 posts on a page
	}

	ip := &IndexPage{}
	ip.CurrentPage = pageNum
	ip.Posts, err = DBGetPosts(10, offset)
	if err != nil {
		log.Printf("Error during db query for index page: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	iph.Template.Execute(w, ip)
}
