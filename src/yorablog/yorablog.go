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

	http.Handle("/login/", LoginHandler)
	http.Handle("/createuser/", CreateUserHandler)
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(BaseStaticPath))))

	http.Handle("/", IndexHandler)
	http.Handle("/post/", PostHandler)
	http.Handle("/edit/", EditHandler)
	http.Handle("/create/", CreateHandler)

	log.Printf("Listen on %v.\n", BaseServeAddr)

	http.ListenAndServe(BaseServeAddr, nil)
}
