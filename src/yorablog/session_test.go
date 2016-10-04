package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type handlerT struct {
	called bool
	req    *http.Request
}

func (h *handlerT) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	h.req = r
}

func TestSessionRequiredNoCookie(t *testing.T) {
	db := &tDB{}
	parent := &handlerT{called: false}
	h := SessionRequired(db, db, parent)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if !parent.called {
		t.Errorf("Parent is not called.\n")
	}

	_, ok := SessionFromContext(parent.req.Context())
	if !ok {
		t.Errorf("Session is not created.\n")
	}
}

func TestSessionRequiredWithCookie(t *testing.T) {
	db := &tDB{}
	parent := &handlerT{called: false}
	h := SessionRequired(db, db, parent)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)
	w := httptest.NewRecorder()

	c := &http.Cookie{}
	c.Name = "SessionID"
	c.Value = "01"

	req.AddCookie(c)

	h.ServeHTTP(w, req)

	if !parent.called {
		t.Errorf("Parent is not called.\n")
	}

	_, ok := SessionFromContext(parent.req.Context())
	if !ok {
		t.Errorf("Session is not created.\n")
	}
}
