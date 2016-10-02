package main

import (
	"log"
	"net/http"
	"yoradb"
)

func main() {

	err := InitEnv()
	if err != nil {
		log.Panicln(err)
	}

	InitURLPatterns() // panic if regexps not compile

	db, err := yoradb.InitDB(BaseDSN)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	pr := db
	sr := db
	ur := db
	rpr := db

	LoginHandler := InitLoginPageHandler(sr, ur, BaseTemplatesPath)
	LogoutHandler := InitLogoutPageHandler(sr)
	CreateUserHandler := InitCreateUserPageHandler(sr, ur, BaseTemplatesPath)
	ForgotPasswordHandler := InitForgotPasswordPageHandler(ur, rpr, BaseTemplatesPath)
	RestorePasswordHandler := InitRestorePasswordPageHandler(rpr, BaseTemplatesPath)

	IndexHandler := InitIndexPageHandler(pr, BaseTemplatesPath)
	PostHandler := InitPostPageHandler(pr, BaseTemplatesPath)
	EditHandler := InitEditPageHandler(pr, BaseTemplatesPath)
	CreateHandler := InitCreatePageHandler(pr, BaseTemplatesPath)

	ErrorTemplate = InitErrorTemplate(BaseTemplatesPath)

	http.Handle("/login/", SessionRequired(sr, ur, LoginHandler))
	http.Handle("/logout/", SessionRequired(sr, ur, LogoutHandler))
	http.Handle("/createuser/", SessionRequired(sr, ur, CreateUserHandler))
	http.Handle("/forgotpassword/", SessionRequired(sr, ur, ForgotPasswordHandler))
	http.Handle("/restorepassword/", SessionRequired(sr, ur, RestorePasswordHandler))

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(BaseStaticPath))))
	http.Handle("/photos/", http.StripPrefix("/photos/",
		http.FileServer(http.Dir(BasePhotosPath))))

	http.Handle("/", SessionRequired(sr, ur, IndexHandler))
	http.Handle("/post/", SessionRequired(sr, ur, PostHandler))
	http.Handle("/edit/", SessionRequired(sr, ur, LoginRequired(sr, ur, EditHandler)))
	http.Handle("/create/", SessionRequired(sr, ur, LoginRequired(sr, ur, CreateHandler)))

	log.Printf("Listen on %v.\n", BaseServeAddr)

	http.ListenAndServe(BaseServeAddr, nil)
}
