package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostPageServeHTTP(t *testing.T) {
	db := &tDB{}
	h := InitPostPageHandler(db, "/home/valya/myprogs/yorablog/templates")
	req := httptest.NewRequest(http.MethodGet, "http://localhost/post/3", nil)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "User", sessions["01"])
	req = req.WithContext(ctx)

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
