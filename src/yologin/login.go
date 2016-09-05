package yologin

import (
	"log"
	"net/http"
	"path/filepath"
	"yotemplate"
)

// LoginPage data for login page
type LoginPage struct {
	Email    string
	Password string
}

// LoginPageHandler - handler for login pages
type LoginPageHandler struct {
	LoginTemplates *yotemplate.YoTemplate
}

// InitLoginPageHandler creates and inits login page handler
func InitLoginPageHandler(templatesPath string) *LoginPageHandler {
	loginTemplatePath := filepath.Join(templatesPath, "login.html")
	loginTemplates, err := yotemplate.InitYoTemplate(loginTemplatePath)
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

	lph.LoginTemplates.Execute(w, lp)
}
