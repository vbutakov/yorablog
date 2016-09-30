package main

import (
	"errors"
	"time"
	"yoradb"
)

type tDB struct {
}

func (db *tDB) CreatePost(post *yoradb.Post, userID int) (int, error) {
	return 10, nil
}

func (db *tDB) GetPostByID(id int) (*yoradb.Post, error) {
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

func (db *tDB) UpdatePost(post *yoradb.Post) error {
	return nil
}

func (db *tDB) GetPosts(num, offset int) ([]yoradb.Post, error) {
	return nil, nil
}

func (db *tDB) DBGetUserBySessionID(sessionID string) (*yoradb.User, error) {

	if user, ok := sessions[sessionID]; ok {
		return user, nil
	}

	return nil, errors.New("User not found")
}

func (db *tDB) DBSessionValid(sessionID string) bool {
	return true
}

func (db *tDB) DBInsertNewSession(sessionID string, expires time.Time) error {
	return nil
}

func (db *tDB) DBCreateUser(name, email, password string) (int, error) {
	return 1, nil
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
	return "a1@b.c", nil
}

func (db *tDB) DBUpdatePasswordByRestoreToken(token, email, password string) error {
	return nil
}

func (db *tDB) Close() error {
	return nil
}
