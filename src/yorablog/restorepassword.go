package main

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"yotemplate"
)

// RestorePasswordPage is a struct for data on the restor password page
type RestorePasswordPage struct {
	Email           string
	Password        string
	PasswordConfirm string
	Message         string
	ErrorMessage    string
}

// RestorePasswordPageHandler is a handler for page processing
type RestorePasswordPageHandler struct {
	Template *yotemplate.Template
}

// InitRestorePasswordPageHandler initialize RestorePasswordPageHandler struct
func InitRestorePasswordPageHandler(templatesPath string) *RestorePasswordPageHandler {
	templatePath := filepath.Join(templatesPath, "restorepassword.html")
	template, err := yotemplate.InitTemplate(templatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Index page template is initialized.")

	return &RestorePasswordPageHandler{Template: template}
}

// RestorePasswordPageHandler - handler for index page
func (rph RestorePasswordPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	rp := &RestorePasswordPage{}
	token := r.FormValue("token")
	if token == "" {
		rp.Message = "Некорректная ссылка для восстановления пароля. Воспользуйтесь ссылкой из письма."

		w.WriteHeader(http.StatusOK)
		rph.Template.Execute(w, rp)
		return
	}

	email, err := DBGetEmailByRestoreToken(token)
	if err == sql.ErrNoRows {
		rp.Message = "Не корректный токен. Для восстановления пароля воспользуйтесь ссылкой из письма."
	} else if err != nil {
		log.Printf("Error in get email by restore token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		rp.Email = email
	}

	if r.Method == "GET" {
	} else if r.Method == "POST" {
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("passwordconfirm")

		if password != passwordConfirm {
			rp.ErrorMessage = "Подтверждение пароля не совпадает с паролем."
		} else if password == "" {
			rp.ErrorMessage = "Пароль не должен быть пустым."
		} else {
			// password is ok
			err = DBUpdatePasswordByRestoreToken(token, email, password)
			if err != nil {
				log.Printf("Error during update password by restore token: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			rp.Message = "Пароль обновлен."
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	rph.Template.Execute(w, rp)
	return
}
