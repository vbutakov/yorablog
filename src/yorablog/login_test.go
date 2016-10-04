package main

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"yoradb"
)

func TestLoginGetServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitLoginPageHandler(db, db, "/home/valya/myprogs/yorablog/templates")
	req := httptest.NewRequest(http.MethodGet, "http://localhost/login", nil)
	w := httptest.NewRecorder()

	session := &yoradb.Session{ID: "04"}
	ctx := context.WithValue(req.Context(), keySession, session)
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("</html>")) {
		t.Errorf("Page is not complete")
	}
}

func TestLoginPostServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitLoginPageHandler(db, db, "/home/valya/myprogs/yorablog/templates")

	form := &url.Values{}
	form.Add("email", "a1@b.c")
	form.Add("password", "111")

	body := strings.NewReader(form.Encode())

	req := httptest.NewRequest(http.MethodPost, "http://localhost/login/?return=/", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	session := &yoradb.Session{ID: "01"}
	ctx := context.WithValue(req.Context(), keySession, session)
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}
}

func TestWrongLoginPostServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitLoginPageHandler(db, db, "/home/valya/myprogs/yorablog/templates")

	form := &url.Values{}
	form.Add("email", "a1@b.c")
	form.Add("password", "222")

	body := strings.NewReader(form.Encode())

	req := httptest.NewRequest(http.MethodPost, "http://localhost/login/?return=/", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	session := &yoradb.Session{ID: "02"}
	ctx := context.WithValue(req.Context(), keySession, session)
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("</html>")) {
		t.Errorf("Page is not complete.\n")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("errormessage")) {
		t.Errorf("No ErrorMessage.\n")
	}
}

func TestLoginRequiredHasUser(t *testing.T) {
	db := &tDB{}
	parent := &handlerT{called: false}
	h := LoginRequired(db, db, parent)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
	w := httptest.NewRecorder()

	userID := sql.NullInt64{Int64: 1, Valid: true}
	s := &yoradb.Session{ID: "01", UserID: userID}
	ctx := context.WithValue(req.Context(), keySession, s)
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if !parent.called {
		t.Errorf("Parent is not called.\n")
	}

	_, ok := UserFromContext(parent.req.Context())
	if !ok {
		t.Errorf("Session is not created.\n")
	}

}

func TestLoginRequiredNoUser(t *testing.T) {
	db := &tDB{}
	parent := &handlerT{called: false}
	h := LoginRequired(db, db, parent)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
	w := httptest.NewRecorder()

	userID := sql.NullInt64{Int64: 1, Valid: false}
	s := &yoradb.Session{ID: "01", UserID: userID}
	ctx := context.WithValue(req.Context(), keySession, s)
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}
}
