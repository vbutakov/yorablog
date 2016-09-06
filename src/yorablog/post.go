package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"yotemplate"
)

// PostPage is a struct for data on the post page
type PostPage struct {
	Post          *Post
	OGURL         string
	OGType        string
	OGTitle       template.HTML
	OGDescription template.HTML
	OGImage       string
}

// PostPageHandler is a handler for post page processing
type PostPageHandler struct {
	Template *yotemplate.YoTemplate
}

// InitPostPageHandler initialize PostPageHandler struct
func InitPostPageHandler(templatesPath string) *PostPageHandler {
	postTemplatePath := filepath.Join(templatesPath, "post.html")
	postTemplate, err := yotemplate.InitYoTemplate(postTemplatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Post page template is initialized.")

	return &PostPageHandler{Template: postTemplate}
}

// PostPageHandler - handler for post page
func (pph PostPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	res := PostURLPattern.FindStringSubmatch(r.URL.Path)
	if res == nil {
		http.NotFound(w, r)
		return
	}

	postID, err := strconv.Atoi(res[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pp := &PostPage{}

	pp.Post, err = DBGetPostByID(postID)
	if err != nil {
		log.Printf("Error during db query for post page: %v\n", err)
		http.NotFound(w, r)
		return
	}

	pp.OGDescription = pp.Post.Description
	pp.OGImage = pp.Post.ImageURL
	pp.OGTitle = pp.Post.Title
	pp.OGType = "website"
	pp.OGURL = "http://" + r.Host + r.URL.String()

	w.WriteHeader(http.StatusOK)

	pph.Template.Execute(w, pp)
}
