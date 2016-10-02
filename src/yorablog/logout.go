package main

import (
	"log"
	"net/http"
	"yoradb"
)

// LogoutPageHandler - handler for logout pages
type LogoutPageHandler struct {
	sr yoradb.SessionRepository
}

// InitLogoutPageHandler initialize CreatePageHandler struct
func InitLogoutPageHandler(sr yoradb.SessionRepository) *LogoutPageHandler {
	return &LogoutPageHandler{sr: sr}
}

func (h LogoutPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, ok := SessionFromContext(r.Context())
	if !ok {
		log.Printf("Error! Cannot logout without session.\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := h.sr.LogoutFromSession(session)
	if err != nil {
		log.Printf("Error during logout: %v\n", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
