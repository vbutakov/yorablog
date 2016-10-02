package main

import (
	"log"
	"net/http"
	"path/filepath"
	"yoradb"
	"yotemplate"
)

// CreateUserPage data for create user page
type CreateUserPage struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
	ErrorMessage    string
	UserName        string
}

// CreateUserPageHandler - handler for create user pages
type CreateUserPageHandler struct {
	template *yotemplate.Template
	sr       yoradb.SessionRepository
	ur       yoradb.UserRepository
}

// InitCreateUserPageHandler creates and inits login page handler
func InitCreateUserPageHandler(sr yoradb.SessionRepository, ur yoradb.UserRepository, templatesPath string) *CreateUserPageHandler {

	pathes := make([]string, 3)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "createuser.gohtml")
	pathes[2] = filepath.Join(templatesPath, "empty_og.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("CreateUser page templates are initialized.")

	return &CreateUserPageHandler{template: templ, sr: sr, ur: ur}
}

// CreateUserHandle - handler for login page
func (h CreateUserPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
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
			h.template.Execute(w, cup)
		} else {
			// create new user
			userID, err := h.ur.CreateUser(cup.Name, cup.Email, cup.Password)
			if err != nil {
				cup.ErrorMessage = err.Error()

				log.Printf("Error during user create: %v\n", err)

				w.WriteHeader(http.StatusOK)
				h.template.Execute(w, cup)

				return
			}

			session, ok := SessionFromContext(r.Context())
			if !ok {
				log.Printf("Error during user create: cannot get session from context.\n")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			session.UserID.Int64 = userID
			err = h.sr.UpdateSession(session)
			if err != nil {
				cup.ErrorMessage = err.Error()

				log.Printf("Error during user create: %v\n", err)

				w.WriteHeader(http.StatusOK)
				h.template.Execute(w, cup)

				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	} else if r.Method == http.MethodGet {
		cup := &CreateUserPage{}
		cup.Name = r.FormValue("name")
		cup.Email = r.FormValue("email")

		w.WriteHeader(http.StatusOK)

		h.template.Execute(w, cup)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
