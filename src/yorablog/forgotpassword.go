package main

import (
	"log"
	"net/http"
	"path/filepath"
	"yotemplate"
)

// ForgotPasswordPage is a struct for data on the forgot password page
type ForgotPasswordPage struct {
	UserEmail    string
	ErrorMessage string
	Message      string
}

// ForgotPasswordPageHandler is a handler for page processing
type ForgotPasswordPageHandler struct {
	Template *yotemplate.Template
}

// InitForgotPasswordPageHandler initialize ForgotPasswordPageHandler struct
func InitForgotPasswordPageHandler(templatesPath string) *ForgotPasswordPageHandler {
	templatePath := filepath.Join(templatesPath, "forgotpassword.html")
	templ, err := yotemplate.InitTemplate(templatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Forgot password page template is initialized.")

	return &ForgotPasswordPageHandler{Template: templ}
}

// ForgotPasswordPageHandler - handler for web page
func (h ForgotPasswordPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		data := &ForgotPasswordPage{}
		w.WriteHeader(http.StatusOK)
		h.Template.Execute(w, data)
		return
	} else if r.Method == "POST" {
		email := r.FormValue("email")
		data := &ForgotPasswordPage{
			UserEmail: email,
		}

		if !DBEmailExist(email) {
			data.ErrorMessage = "Пользователя с таким email-ом не зарегистрировано."
			w.WriteHeader(http.StatusOK)
			h.Template.Execute(w, data)
			return
		}

		id, err := DBCreateRestorePasswordID(email)
		if err != nil {
			log.Printf("Error during password restore: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
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
		h.Template.Execute(w, data)
		return
	}
}
