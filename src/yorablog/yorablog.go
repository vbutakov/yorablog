package main

import (
	"log"
	"net/http"
)

func main() {

	err := InitEnv()
	if err != nil {
		log.Panicln(err)
	}

	InitURLPatterns() // panic if regexps not compile

	db, err := InitDB()
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()
	DBConnection = db

	testHandler := InitTestPostHandler(db, BaseTemplatesPath)

	LoginHandler := InitLoginPageHandler(BaseTemplatesPath)
	LogoutHandler := &LogoutPageHandler{}
	CreateUserHandler := InitCreateUserPageHandler(BaseTemplatesPath)
	ForgotPasswordHandler := InitForgotPasswordPageHandler(BaseTemplatesPath)
	RestorePasswordHandler := InitRestorePasswordPageHandler(BaseTemplatesPath)

	IndexHandler := InitIndexPageHandler(BaseTemplatesPath)
	PostHandler := InitPostPageHandler(BaseTemplatesPath)
	EditHandler := InitEditPageHandler(BaseTemplatesPath)
	CreateHandler := InitCreatePageHandler(BaseTemplatesPath)

	ErrorTemplate = InitErrorTemplate(BaseTemplatesPath)

	http.Handle("/login/", SessionRequired(LoginHandler))
	http.Handle("/logout/", SessionRequired(LogoutHandler))
	http.Handle("/createuser/", SessionRequired(CreateUserHandler))
	http.Handle("/forgotpassword/", SessionRequired(ForgotPasswordHandler))
	http.Handle("/restorepassword/", SessionRequired(RestorePasswordHandler))

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(BaseStaticPath))))
	http.Handle("/photos/", http.StripPrefix("/photos/",
		http.FileServer(http.Dir(BasePhotosPath))))

	http.Handle("/", SessionRequired(IndexHandler))
	http.Handle("/post/", SessionRequired(PostHandler))
	http.Handle("/edit/", SessionRequired(LoginRequired(EditHandler)))
	http.Handle("/create/", SessionRequired(LoginRequired(CreateHandler)))

	http.Handle("/test_post/", SessionRequired(testHandler))

	log.Printf("Listen on %v.\n", BaseServeAddr)

	http.ListenAndServe(BaseServeAddr, nil)
}
