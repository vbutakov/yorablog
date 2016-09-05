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

	LoginHandler := yologin.InitLoginPageHandler(BaseTemplatesPath)
	CreateUserHandler := yologin.InitCreateUserPageHandler(BaseTemplatesPath)

	IndexHandler := InitIndexPageHandler(BaseTemplatesPath)

	http.Handle("/login/", LoginHandler)
	http.Handle("/createuser/", CreateUserHandler)
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(BaseStaticPath))))

	http.Handle("/", IndexHandler)

	log.Printf("Listen on %v.\n", BaseServeAddr)

	http.ListenAndServe(BaseServeAddr, nil)
}
