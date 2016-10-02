package main

import (
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

func (db *tDB) CreateUser(name, email, password string) (int, error) {
	return 1, nil
}

func (db *tDB) LoginUser(email, password string) (int, error) {
	return 1, nil
}

func (db *tDB) GetUserByID(id int) (*User, error) {

}

func (db *tDB) GetUserByEmail(email string) (*User, error) {

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
