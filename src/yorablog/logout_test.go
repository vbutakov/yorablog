package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"yoradb"
)

func TestLogoutGetServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitLogoutPageHandler(db)
	req := httptest.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	w := httptest.NewRecorder()

	session := &yoradb.Session{ID: "01"}
	ctx := context.WithValue(req.Context(), keySession, session)
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}
}
