package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"yotemplate"
)

// IndexPage is a struct for data on the index page
type IndexPage struct {
	CurrentPageURL string
	PrevPageURL    string
	NextPageURL    string
	UserName       string
	Posts          []Post
}

// IndexPageHandler is a handler for page processing
type IndexPageHandler struct {
	Template *yotemplate.Template
}

// InitIndexPageHandler initialize IndexPageHandler struct
func InitIndexPageHandler(templatesPath string) *IndexPageHandler {
	indexTemplatePath := filepath.Join(templatesPath, "index.html")
	indexTemplate, err := yotemplate.InitTemplate(indexTemplatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Index page template is initialized.")

	return &IndexPageHandler{Template: indexTemplate}
}

// IndexPageHandle - handler for index page
func (iph IndexPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	res := IndexURLPattern.FindStringSubmatch(r.URL.Path)
	if res == nil {
		http.NotFound(w, r)
		return
	}

	var offset int
	var pageNum int
	var err error

	ip := &IndexPage{}

	cookie, err := r.Cookie("SessionID")
	if err != nil {
		log.Printf("Error during cookie read from on index page: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sessionID := cookie.Value
	user, err := DBGetUserBySessionID(DBConnection, sessionID)
	if err == nil {
		ip.UserName = user.Name
	}

	// calculate offset for db query
	pageNumStr := res[1]
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

	// create links for prev/next buttons
	ip.CurrentPageURL = r.URL.String()
	if offset == 0 {
		ip.PrevPageURL = "/2"
		ip.NextPageURL = "/"
	} else {
		ip.PrevPageURL = fmt.Sprintf("/%v", pageNum+1)
		if pageNum > 2 {
			ip.NextPageURL = fmt.Sprintf("/%v", pageNum-1)
		} else {
			ip.NextPageURL = "/"
		}
	}

	ip.Posts, err = DBGetPosts(10, offset)
	if err != nil {
		log.Printf("Error during db query for index page: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	iph.Template.Execute(w, ip)
}
