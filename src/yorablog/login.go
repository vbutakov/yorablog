package main

import (
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"yotemplate"
)

// LoginPage data for login page
type LoginPage struct {
	Email        string
	Password     string
	ErrorMessage string
	URLQuery     string
}

// LoginPageHandler - handler for login pages
type LoginPageHandler struct {
	LoginTemplates *yotemplate.YoTemplate
}

// LoginRequiredHandler structure for checking if user login required
type LoginRequiredHandler struct {
	parent http.Handler
}

// LoginRequired initialize LoginRequiredHandler
func LoginRequired(parent http.Handler) LoginRequiredHandler {
	return LoginRequiredHandler{parent}
}

func (h LoginRequiredHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		log.Printf("Error! Cannot check if login required without session: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !DBUserIsLogedIn(cookie.Value) {
		val := make(url.Values)
		val.Add("return", r.URL.String())
		redirectURL := "/login/?" + val.Encode()
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		h.parent.ServeHTTP(w, r)
	}
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

	if r.Method == "POST" {
		lp := &LoginPage{URLQuery: r.URL.RawQuery}
		lp.Email = r.FormValue("email")
		lp.Password = r.FormValue("password")

		userID, err := DBLoginUser(lp.Email, lp.Password)
		if err == ErrLoginFailed {
			lp.ErrorMessage = "Не угадали email и пароль"
			w.WriteHeader(http.StatusOK)
			lph.LoginTemplates.Execute(w, lp)
			return
		} else if err != nil {
			log.Printf("Error during user login: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie, err := r.Cookie("SessionID")
		if err != nil {
			log.Printf("Error during user login: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sessionID := cookie.Value
		err = DBUpdateSessionWithUserID(sessionID, userID)
		if err != nil {
			log.Printf("Error during user login: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		redirectURL := r.URL.Query().Get("return")
		if redirectURL == "" {
			redirectURL = "/"
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)

	} else if r.Method == "GET" {
		lp := &LoginPage{URLQuery: r.URL.RawQuery}
		lp.Email = r.FormValue("email")

		w.WriteHeader(http.StatusOK)
		lph.LoginTemplates.Execute(w, lp)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
