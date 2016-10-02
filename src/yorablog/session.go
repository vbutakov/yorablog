package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"yoradb"
)

type key string

const (
	keySession = "Session"
	keyUser    = "User"
)

// SessionHandler is the handler for session creation
type SessionHandler struct {
	parent http.Handler
	sr     yoradb.SessionRepository
	ur     yoradb.UserRepository
}

// SessionRequired initialize session handler
func SessionRequired(sr yoradb.SessionRepository,
	ur yoradb.UserRepository,
	parent http.Handler) SessionHandler {
	return SessionHandler{parent: parent, sr: sr, ur: ur}
}

func (h SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var needSetCookie = false
	var sessionID string
	var session *yoradb.Session

	cookie, err := r.Cookie("SessionID")
	if err != nil {
		needSetCookie = true
		sessionID = CreateSessionID()
	} else {
		sessionID = cookie.Value
		session, err = h.sr.GetSessionByID(sessionID)
		if err != nil {
			needSetCookie = true
			sessionID = CreateSessionID()
		} else {
			if !session.UserID.Valid {
				user, err := h.ur.GetUserByID(session.UserID.Int64)
				if err != nil {
					h.sr.DeleteSession(session)
					needSetCookie = true
					sessionID = CreateSessionID()
				} else {
					uctx := context.WithValue(r.Context(), keyUser, user)
					r = r.WithContext(uctx)
				}
			}
		}
	}

	if needSetCookie {
		expires := time.Now().Add(30 * 24 * time.Hour)

		cookie = &http.Cookie{}
		cookie.Name = "SessionID"
		cookie.Value = sessionID
		cookie.Path = "/"
		cookie.Expires = expires

		// have to reload page because it need cookie
		w.Header().Add("Set-Cookie", cookie.String())

		session = &yoradb.Session{ID: sessionID, Expires: expires}
		h.sr.CreateSession(session)
	}

	// set sessionID into context
	ctx := context.WithValue(r.Context(), keySession, session)
	r = r.WithContext(ctx)

	h.parent.ServeHTTP(w, r)
}

// CreateSessionID create new random session identifier
func CreateSessionID() string {
	nsec := time.Now().UnixNano()
	rand.Seed(nsec)
	r := rand.Int63n(1e9)
	s := fmt.Sprintf("%v_%v", nsec, r)
	return s
}

// SessionFromContext get session from request context
func SessionFromContext(ctx context.Context) (*yoradb.Session, bool) {
	session, ok := ctx.Value(keySession).(*yoradb.Session)

	return session, ok
}

// UserFromContext get session from request context
func UserFromContext(ctx context.Context) (*yoradb.User, bool) {
	session, ok := ctx.Value(keyUser).(*yoradb.User)

	return session, ok
}
