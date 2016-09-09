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

	err = InitDB()
	if err != nil {
		log.Panicln(err)
	}
	defer DBConnection.Close()

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

	http.Handle("/", SessionRequired(IndexHandler))
	http.Handle("/post/", SessionRequired(PostHandler))
	http.Handle("/edit/", SessionRequired(LoginRequired(EditHandler)))
	http.Handle("/create/", SessionRequired(LoginRequired(CreateHandler)))

	log.Printf("Listen on %v.\n", BaseServeAddr)

	http.ListenAndServe(BaseServeAddr, nil)
}
