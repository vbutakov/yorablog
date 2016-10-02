package main

import (
	"bytes"
	"context"
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
	form.Add("email", "User1")
	form.Add("password", "a1@b.c")

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
