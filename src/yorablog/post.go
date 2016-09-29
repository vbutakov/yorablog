package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"yoradb"
	"yotemplate"
)

// PostPage is a struct for data on the post page
type PostPage struct {
	Post          *yoradb.Post
	OGURL         string
	OGType        string
	OGTitle       template.HTML
	OGDescription template.HTML
	OGImage       string
	UserName      string
	ErrorMessage  string
}

// PostPageHandler is a handler for post page processing
type PostPageHandler struct {
	template *yotemplate.Template
	db       yoradb.PostRepository
}

// InitPostPageHandler initialize PostPageHandler struct
func InitPostPageHandler(db yoradb.PostRepository, templatesPath string) *PostPageHandler {

	pathes := make([]string, 3)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "post.gohtml")
	pathes[2] = filepath.Join(templatesPath, "post_og.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Post page template is initialized.")

	return &PostPageHandler{template: templ, db: db}
}

// PostPageHandler - handler for post page
func (h PostPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	res := PostURLPattern.FindStringSubmatch(r.URL.Path)
	if res == nil {
		http.NotFound(w, r)
		return
	}

	pp := &PostPage{}

	user, ok := r.Context().Value("User").(*yoradb.User)
	if ok {
		pp.UserName = user.Name
	}

	postID, err := strconv.Atoi(res[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pp.Post, err = h.db.DBGetPostByID(postID)
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
