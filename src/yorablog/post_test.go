package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"yoradb"
)

type tDB struct {
}

func (db *tDB) DBInsertPost(post *yoradb.Post, userID int) (int, error) {
	return 0, nil
}

func (db *tDB) DBGetPostByID(id int) (*yoradb.Post, error) {
	p := &yoradb.Post{
		ID:          id,
		Title:       "Тестовый заголовок",
		Description: "Тестовое описание",
		ImageURL:    "http://example.com/image.jpg",
		Annotation:  "Тестовая аннотация",
		Text:        "Тестовый текст",
		Author:      "Тестовый автор",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return p, nil
}

func (db *tDB) DBUpdatePost(post *yoradb.Post) error {
	return nil
}

func (db *tDB) DBGetPosts(num, offset int) ([]yoradb.Post, error) {
	return nil, nil
}

func (db *tDB) DBGetUserBySessionID(sessionID string) (*yoradb.User, error) {
	user := &yoradb.User{}
	user.ID = 10
	user.Name = "Begemotina"
	user.Email = "vozhdyara@gmail.com"

	return user, nil
}

func (db *tDB) DBSessionValid(sessionID string) bool {
	return true
}

func (db *tDB) DBInsertNewSession(sessionID string, expires time.Time) error {
	return nil
}

func (db *tDB) DBUserIsLogedIn(sessionID string) bool {
	return true
}

func (db *tDB) DBCreateUser(name, email, password string) (int, error) {
	return 10, nil
}

func (db *tDB) DBUpdateSessionWithUserID(sessionID string, userID int) error {
	return nil
}

func (db *tDB) DBLoginUser(email, password string) (int, error) {
	return 10, nil
}

func (db *tDB) DBLogoutUserFromSession(sessionID string) error {
	return nil
}

func (db *tDB) DBEmailExist(email string) bool {
	return true
}

func (db *tDB) DBCreateRestorePasswordID(email, token string) (string, error) {
	return "", nil
}

func (db *tDB) DBGetEmailByRestoreToken(token string) (string, error) {
	return "", nil
}

func (db *tDB) DBUpdatePasswordByRestoreToken(token, email, password string) error {
	return nil
}

func (db *tDB) Close() error {
	return nil
}

func TestPostPageServeHTTP(t *testing.T) {
	db := &tDB{}
	InitURLPatterns()
	h := InitPostPageHandler(db, "/home/valya/myprogs/yorablog/templates")
	req := httptest.NewRequest("GET", "http://localhost/post/3", nil)
	w := httptest.NewRecorder()

	c := &http.Cookie{}
	c.Name = "SessionID"
	c.Value = "1234567890"

	req.Header.Set("Cookie", c.String())

	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Wrong error code: %v\n", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("</html>")) {
		t.Errorf("Page is not complete")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
