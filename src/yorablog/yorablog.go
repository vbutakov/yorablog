package main

import (
	"log"
	"net/http"
	"yologin"
)

func main() {

	err := InitEnv()
	if err != nil {
		log.Panicln(err)
	}

	InitURLPatterns() // panic if regexps not compile

	err = InitDB()
	if err != nil {
		log.Panicln(err)
	}
	defer DBConnection.Close()

	LoginHandler := yologin.InitLoginPageHandler(BaseTemplatesPath)
	CreateUserHandler := yologin.InitCreateUserPageHandler(BaseTemplatesPath)

	IndexHandler := InitIndexPageHandler(BaseTemplatesPath)
	PostHandler := InitPostPageHandler(BaseTemplatesPath)
	EditHandler := InitEditPageHandler(BaseTemplatesPath)
	CreateHandler := InitCreatePageHandler(BaseTemplatesPath)

	http.Handle("/login/", SessionRequired(LoginHandler))
	http.Handle("/createuser/", SessionRequired(CreateUserHandler))
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(BaseStaticPath))))

	http.Handle("/", SessionRequired(IndexHandler))
	http.Handle("/post/", SessionRequired(PostHandler))
	http.Handle("/edit/", SessionRequired(EditHandler))
	http.Handle("/create/", SessionRequired(CreateHandler))

	log.Printf("Listen on %v.\n", BaseServeAddr)

	http.ListenAndServe(BaseServeAddr, nil)
}
