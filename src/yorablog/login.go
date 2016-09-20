package main

import (
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"yoradb"
	"yotemplate"
)

// LoginPage data for login page
type LoginPage struct {
	Email        string
	Password     string
	ErrorMessage string
	URLQuery     string
	UserName     string
}

// LoginPageHandler - handler for login pages
type LoginPageHandler struct {
	template *yotemplate.Template
	db       yoradb.DB
}

// LoginRequiredHandler structure for checking if user login required
type LoginRequiredHandler struct {
	parent http.Handler
	db     yoradb.DB
}

// LoginRequired initialize LoginRequiredHandler
func LoginRequired(db yoradb.DB, parent http.Handler) LoginRequiredHandler {
	return LoginRequiredHandler{parent: parent, db: db}
}

func (h LoginRequiredHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		log.Printf("Error! Cannot check if login required without session: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !h.db.DBUserIsLogedIn(cookie.Value) {
		val := make(url.Values)
		val.Add("return", r.URL.String())
		redirectURL := "/login/?" + val.Encode()
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		h.parent.ServeHTTP(w, r)
	}
}

// InitLoginPageHandler creates and inits login page handler
func InitLoginPageHandler(db yoradb.DB, templatesPath string) *LoginPageHandler {

	pathes := make([]string, 3)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "login.gohtml")
	pathes[2] = filepath.Join(templatesPath, "empty_og.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Login page templates are initialized.")

	return &LoginPageHandler{template: templ, db: db}
}

// LoginHandle - handler for login page
func (h LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		lp := &LoginPage{URLQuery: r.URL.RawQuery}
		lp.Email = r.FormValue("email")
		lp.Password = r.FormValue("password")

		userID, err := h.db.DBLoginUser(lp.Email, lp.Password)
		if err == yoradb.ErrLoginFailed {

			lp.ErrorMessage = "Не угадали email и пароль"

			w.WriteHeader(http.StatusOK)
			h.template.Execute(w, lp)
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
		err = h.db.DBUpdateSessionWithUserID(sessionID, userID)
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
		h.template.Execute(w, lp)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
