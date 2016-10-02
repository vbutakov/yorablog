package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestForgotPasswordServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitForgotPasswordPageHandler(db, db, "/home/valya/myprogs/yorablog/templates")
	req := httptest.NewRequest(http.MethodGet, "http://localhost/forgotpassword/", nil)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "User", sessions["04"])
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("</html>")) {
		t.Errorf("Page is not complete")
	}
}
