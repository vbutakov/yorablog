package yologin

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// CreateUserPage data for create user page
type CreateUserPage struct {
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
}

// CreateUserPageHandler - handler for create user pages
type CreateUserPageHandler struct {
	CreateUserTemplates *template.Template
}

// InitCreateUserPageHandler creates and inits login page handler
func InitCreateUserPageHandler() *CreateUserPageHandler {
	createUserTemplatePath := filepath.Join("templates", "createuser.html")
	createUserTemplates, err := template.ParseFiles(createUserTemplatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("CreateUser page templates are initialized.")

	return &CreateUserPageHandler{CreateUserTemplates: createUserTemplates}
}

// CreateUserHandle - handler for login page
func (cuph CreateUserPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

	}
	cup := &CreateUserPage{}
	cup.Name = r.FormValue("name")
	cup.Email = r.FormValue("email")

	w.WriteHeader(http.StatusOK)

	cuph.CreateUserTemplates.ExecuteTemplate(w, "createuser.html", cup)
}