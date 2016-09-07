package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// SessionHandler is the handler for session creation
type SessionHandler struct {
	parent http.Handler
}

// SessionRequired initialize session handler
func SessionRequired(parent http.Handler) SessionHandler {
	return SessionHandler{parent}
}

func (h SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var needSetCookie = false
	var sessionID string

	cookie, err := r.Cookie("SessionID")
	if err != nil {
		needSetCookie = true
		sessionID = createSessionID()
	} else {
		sessionID = cookie.Value
		if !DBSessionValid(sessionID) {
			needSetCookie = true
			sessionID = createSessionID()
		}
	}

	if needSetCookie {
		expires := time.Now().Add(30 * 24 * time.Hour)
		_ = DBInsertNewSession(sessionID, expires)

		cookie = &http.Cookie{}
		cookie.Name = "SessionID"
		cookie.Value = sessionID
		cookie.Path = "/"
		cookie.Expires = expires
		w.Header().Add("Set-Cookie", cookie.String())
	}

	h.parent.ServeHTTP(w, r)
}

func createSessionID() string {
	nsec := time.Now().UnixNano()
	rand.Seed(nsec)
	r := rand.Int63n(1e9)
	s := fmt.Sprintf("%v_%v", nsec, r)
	return s
}
