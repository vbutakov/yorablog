package main

import (
	"context"
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
	sr       yoradb.SessionRepository
	ur       yoradb.UserRepository
}

// LoginRequiredHandler structure for checking if user login required
type LoginRequiredHandler struct {
	parent http.Handler
	sr     yoradb.SessionRepository
	ur     yoradb.UserRepository
}

// LoginRequired initialize LoginRequiredHandler
func LoginRequired(sr yoradb.SessionRepository, ur yoradb.UserRepository, parent http.Handler) LoginRequiredHandler {
	return LoginRequiredHandler{parent: parent, sr: sr, ur: ur}
}

func (h LoginRequiredHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	session, ok := SessionFromContext(r.Context())
	if !ok {
		log.Printf("Error! Cannot check if login required without session.\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !session.UserID.Valid {
		val := make(url.Values)
		val.Add("return", r.URL.String())
		redirectURL := "/login/?" + val.Encode()
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		user, err := h.ur.GetUserByID(session.UserID.Int64)
		if err != nil {
			log.Printf("Error! Cannot read user from db on login page: %v.\n", err)
			val := make(url.Values)
			val.Add("return", r.URL.String())
			redirectURL := "/login/?" + val.Encode()
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		} else {
			uctx := context.WithValue(r.Context(), keyUser, user)
			r = r.WithContext(uctx)

			h.parent.ServeHTTP(w, r)
		}
	}
}

// InitLoginPageHandler creates and inits login page handler
func InitLoginPageHandler(sr yoradb.SessionRepository, ur yoradb.UserRepository, templatesPath string) *LoginPageHandler {

	pathes := make([]string, 3)
	pathes[0] = filepath.Join(templatesPath, "layout.gohtml")
	pathes[1] = filepath.Join(templatesPath, "login.gohtml")
	pathes[2] = filepath.Join(templatesPath, "empty_og.gohtml")

	templ, err := yotemplate.InitTemplate(pathes...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Login page templates are initialized.")

	return &LoginPageHandler{template: templ, sr: sr, ur: ur}
}

// LoginHandle - handler for login page
func (h LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		lp := &LoginPage{URLQuery: r.URL.RawQuery}
		lp.Email = r.FormValue("email")
		lp.Password = r.FormValue("password")

		userID, err := h.ur.LoginUser(lp.Email, lp.Password)
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

		session, ok := SessionFromContext(r.Context())
		if !ok {
			log.Printf("Error during user login: cannot read session.\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session.UserID.Int64 = userID
		err = h.sr.UpdateSession(session)
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

	} else if r.Method == http.MethodGet {
		lp := &LoginPage{URLQuery: r.URL.RawQuery}
		lp.Email = r.FormValue("email")

		w.WriteHeader(http.StatusOK)
		h.template.Execute(w, lp)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
