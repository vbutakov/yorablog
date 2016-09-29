package main

import (
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

	c := &http.Cookie{}
	c.Name = "SessionID"
	c.Value = "04"

	req.Header.Set("Cookie", c.String())

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

	c := &http.Cookie{}
	c.Name = "SessionID"
	c.Value = "01"

	req.Header.Set("Cookie", c.String())

	h.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}
}
