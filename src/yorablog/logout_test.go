package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoutGetServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitLogoutPageHandler(db)
	req := httptest.NewRequest(http.MethodGet, "http://localhost/logout", nil)
	w := httptest.NewRecorder()

	c := &http.Cookie{}
	c.Name = "SessionID"
	c.Value = "01"

	req.Header.Set("Cookie", c.String())

	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}
}
