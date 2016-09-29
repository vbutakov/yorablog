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

	pr, ok := db.(yoradb.PostRepository)
	if !ok {
		log.Panicln("Error ehile converting db to PostRepository")
	}

	LoginHandler := InitLoginPageHandler(db, BaseTemplatesPath)
	LogoutHandler := InitLogoutPageHandler(db)
	CreateUserHandler := InitCreateUserPageHandler(db, BaseTemplatesPath)
	ForgotPasswordHandler := InitForgotPasswordPageHandler(db, BaseTemplatesPath)
	RestorePasswordHandler := InitRestorePasswordPageHandler(db, BaseTemplatesPath)

	IndexHandler := InitIndexPageHandler(pr, BaseTemplatesPath)
	PostHandler := InitPostPageHandler(pr, BaseTemplatesPath)
	EditHandler := InitEditPageHandler(pr, BaseTemplatesPath)
	CreateHandler := InitCreatePageHandler(pr, BaseTemplatesPath)

	ErrorTemplate = InitErrorTemplate(BaseTemplatesPath)

	http.Handle("/login/", SessionRequired(db, LoginHandler))
	http.Handle("/logout/", SessionRequired(db, LogoutHandler))
	http.Handle("/createuser/", SessionRequired(db, CreateUserHandler))
	http.Handle("/forgotpassword/", SessionRequired(db, ForgotPasswordHandler))
	http.Handle("/restorepassword/", SessionRequired(db, RestorePasswordHandler))

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(BaseStaticPath))))
	http.Handle("/photos/", http.StripPrefix("/photos/",
		http.FileServer(http.Dir(BasePhotosPath))))

	http.Handle("/", SessionRequired(db, IndexHandler))
	http.Handle("/post/", SessionRequired(db, PostHandler))
	http.Handle("/edit/", SessionRequired(db, LoginRequired(db, EditHandler)))
	http.Handle("/create/", SessionRequired(db, LoginRequired(db, CreateHandler)))

	log.Printf("Listen on %v.\n", BaseServeAddr)

	http.ListenAndServe(BaseServeAddr, nil)
}
