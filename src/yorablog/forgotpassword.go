package main

import (
	"log"
	"net/http"
	"path/filepath"
	"yoradb"
	"yotemplate"
)

// ForgotPasswordPage is a struct for data on the forgot password page
type ForgotPasswordPage struct {
	UserEmail    string
	ErrorMessage string
	Message      string
	UserName     string
}

// ForgotPasswordPageHandler is a handler for page processing
type ForgotPasswordPageHandler struct {
	template *yotemplate.Template
	ur       yoradb.UserRepository
	rpr      yoradb.RestorePasswordRepository
}

// InitForgotPasswordPageHandler initialize ForgotPasswordPageHandler struct
func InitForgotPasswordPageHandler(ur yoradb.UserRepository,
	rpr yoradb.RestorePasswordRepository,
	templatesPath string) *ForgotPasswordPageHandler {

	pathes := make([]string, 3)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "forgotpassword.gohtml")
	pathes[2] = filepath.Join(templatesPath, "empty_og.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Forgot password page template is initialized.")

	return &ForgotPasswordPageHandler{template: templ, ur: ur, rpr: rpr}
}

// ForgotPasswordPageHandler - handler for web page
func (h ForgotPasswordPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		data := &ForgotPasswordPage{}
		w.WriteHeader(http.StatusOK)
		h.template.Execute(w, data)
		return
	} else if r.Method == http.MethodPost {
		email := r.FormValue("email")
		data := &ForgotPasswordPage{
			UserEmail: email,
		}

		_, err := h.ur.GetUserByEmail(email)
		if err != nil {
			data.ErrorMessage = "Пользователя с таким email-ом не зарегистрировано."
			w.WriteHeader(http.StatusOK)
			h.template.Execute(w, data)
			return
		}

		token := CreateSessionID()
		id, err := h.rpr.CreateRestorePasswordID(email, token)
		if err != nil {
			data.ErrorMessage = err.Error()

			log.Printf("Error during password restore: %v\n", err)

			w.WriteHeader(http.StatusOK)
			h.template.Execute(w, data)

			return
		}

		// send email
		err = SendEmailForPasswordRestore(email, id)
		if err != nil {
			log.Printf("Error during email send: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data.Message = "На email выслано письмо с инструкцией как восстановить пароль."

		w.WriteHeader(http.StatusOK)
		h.template.Execute(w, data)
		return
	}
}
