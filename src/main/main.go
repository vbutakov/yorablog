package main

import (
	"log"
	"net/http"
	"yologin"
)

func main() {

	LoginHandler := yologin.InitLoginPageHandler()

	http.Handle("/login/", LoginHandler)
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))

	log.Println("Listen on port 8080.")

	http.ListenAndServe(":8080", nil)
}
