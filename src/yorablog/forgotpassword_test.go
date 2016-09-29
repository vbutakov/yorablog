package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestForgotPasswordServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitForgotPasswordPageHandler(db, "/home/valya/myprogs/yorablog/templates")
	req := httptest.NewRequest(http.MethodGet, "http://localhost/forgotpassword/", nil)
	w := httptest.NewRecorder()

	c := &http.Cookie{}
	c.Name = "SessionID"
	c.Value = "01"

	req.Header.Set("Cookie", c.String())

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("</html>")) {
		t.Errorf("Page is not complete")
	}
}
