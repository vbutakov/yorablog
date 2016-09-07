package main

import (
	"log"
	"net/http"
)

// LogoutPageHandler - handler for logout pages
type LogoutPageHandler struct {
}

func (h LogoutPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionID")
	if err != nil {
		log.Printf("Error! Cannot check if login required without session: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = DBLogoutUserFromSession(cookie.Value)
	if err != nil {
		log.Printf("Error during logout: %v\n", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
