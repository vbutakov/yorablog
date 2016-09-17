package main

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"yoradb"
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
	template *yotemplate.Template
	db       *yoradb.DB
}

// InitRestorePasswordPageHandler initialize RestorePasswordPageHandler struct
func InitRestorePasswordPageHandler(db *yoradb.DB, templatesPath string) *RestorePasswordPageHandler {

	pathes := make([]string, 2)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "restorepassword.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Restore password page templates are initialized.")

	return &RestorePasswordPageHandler{template: templ, db: db}
}

// RestorePasswordPageHandler - handler for index page
func (h RestorePasswordPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	rp := &RestorePasswordPage{}
	token := r.FormValue("token")
	if token == "" {
		rp.Message = "Некорректная ссылка для восстановления пароля. Воспользуйтесь ссылкой из письма."

		w.WriteHeader(http.StatusOK)
		h.template.Execute(w, rp)
		return
	}

	email, err := h.db.DBGetEmailByRestoreToken(token)
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
			err = h.db.DBUpdatePasswordByRestoreToken(token, email, password)
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
	h.template.Execute(w, rp)
	return
}
