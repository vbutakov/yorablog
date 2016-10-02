package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreatePageHavePermitServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitCreatePageHandler(db, "/home/valya/myprogs/yorablog/templates")

	form := &url.Values{}
	form.Add("title", "Тестовый заголовок")
	form.Add("description", "Тестовое описание")
	form.Add("imageurl", "Тестовая картинка")
	form.Add("annotation", "Тестовая аннотация")
	form.Add("posttext", "Тестовый текст")

	body := strings.NewReader(form.Encode())

	req := httptest.NewRequest(http.MethodPost, "http://localhost/create", body)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), keyUser, sessions["04"])
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}

	loc := w.Header().Get("Location")
	if !strings.Contains(loc, "/post/10") {
		t.Errorf("Wrong redirect: %v\n", loc)
	}
}

func TestCreatePageNotHavePermitServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitCreatePageHandler(db, "/home/valya/myprogs/yorablog/templates")

	form := &url.Values{}
	form.Add("title", "Тестовый заголовок")
	form.Add("description", "Тестовое описание")
	form.Add("imageurl", "Тестовая картинка")
	form.Add("annotation", "Тестовая аннотация")
	form.Add("posttext", "Тестовый текст")

	body := strings.NewReader(form.Encode())

	req := httptest.NewRequest(http.MethodPost, "http://localhost/create", body)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "User", sessions["01"])
	req = req.WithContext(ctx)

	h.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("</html>")) {
		t.Errorf("Page is not complete")
	}

}
