package yologin

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// LoginPage data for login page
type LoginPage struct {
	Email    string
	Password string
}

// LoginPageHandler - handler for login pages
type LoginPageHandler struct {
	LoginTemplates *template.Template
}

// InitLoginPageHandler creates and inits login page handler
func InitLoginPageHandler() *LoginPageHandler {
	loginTemplatePath := filepath.Join("template", "login.html")
	loginTemplates, err := template.ParseFiles(loginTemplatePath)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Login page templates are initialized.")

	return &LoginPageHandler{LoginTemplates: loginTemplates}
}

// LoginHandle - handler for login page
func (lph LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

	}
	lp := &LoginPage{}
	lp.Email = r.FormValue("email")

	w.WriteHeader(http.StatusOK)

	lph.LoginTemplates.ExecuteTemplate(w, "login.html", lp)
}
