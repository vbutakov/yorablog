package main

import (
	"log"
	"net/http"
	"path/filepath"
	"yotemplate"
)

// CreateUserPage data for create user page
type CreateUserPage struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
	ErrorMessage    string
}

// CreateUserPageHandler - handler for create user pages
type CreateUserPageHandler struct {
	CreateUserTemplates *yotemplate.Template
}

// InitCreateUserPageHandler creates and inits login page handler
func InitCreateUserPageHandler(templatesPath string) *CreateUserPageHandler {
	createUserTemplatePath := filepath.Join(templatesPath, "createuser.html")
	createUserTemplates, err := yotemplate.InitTemplate(createUserTemplatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("CreateUser page templates are initialized.")

	return &CreateUserPageHandler{CreateUserTemplates: createUserTemplates}
}

// CreateUserHandle - handler for login page
func (cuph CreateUserPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var sendFormAgain = false
		cup := &CreateUserPage{}
		cup.Name = r.FormValue("name")
		cup.Email = r.FormValue("email")
		cup.Password = r.FormValue("password")
		cup.ConfirmPassword = r.FormValue("password_confirm")

		if cup.Password != cup.ConfirmPassword {
			cup.ErrorMessage = "Пароль не совпадает с подтверждением"
			sendFormAgain = true
		}

		if cup.Password == "" {
			cup.ErrorMessage = "Пароль не должен быть пустым"
			sendFormAgain = true
		}

		if cup.Email == "" {
			cup.ErrorMessage = "Email не должен быть пустым"
			sendFormAgain = true
		}

		if cup.Name == "" {
			cup.ErrorMessage = "Имя не должно быть пустым"
			sendFormAgain = true
		}

		if sendFormAgain {
			w.WriteHeader(http.StatusOK)
			cuph.CreateUserTemplates.Execute(w, cup)
		} else {
			// create new user
			userID, err := DBCreateUser(cup.Name, cup.Email, cup.Password)
			if err != nil {
				log.Printf("Error during user create: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			cookie, err := r.Cookie("SessionID")
			if err != nil {
				log.Printf("Error during user create: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			sessionID := cookie.Value
			err = DBUpdateSessionWithUserID(sessionID, userID)
			if err != nil {
				log.Printf("Error during user create: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	} else if r.Method == "GET" {
		cup := &CreateUserPage{}
		cup.Name = r.FormValue("name")
		cup.Email = r.FormValue("email")

		w.WriteHeader(http.StatusOK)

		cuph.CreateUserTemplates.Execute(w, cup)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
