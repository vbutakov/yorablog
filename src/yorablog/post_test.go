package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostPageServeHTTP(t *testing.T) {
	db := &tDB{}
	temp := InitPostPageHandler(db, "/home/valya/myprogs/yorablog/templates")
	h := SessionRequired(db, temp)
	req := httptest.NewRequest(http.MethodGet, "http://localhost/post/3", nil)
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

	if !bytes.Contains(w.Body.Bytes(), []byte("Тестовый заголовок")) {
		t.Errorf("Заголовок не выведен")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Тестовый текст")) {
		t.Errorf("Текст не выведен")
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Тестовый автор")) {
		t.Errorf("Автор не выведен")
	}
}
