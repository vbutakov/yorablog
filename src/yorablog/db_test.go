package main

import (
	"errors"
	"time"
	"yoradb"
)

type tDB struct {
}

func (db *tDB) CreatePost(post *yoradb.Post, userID int64) (int64, error) {
	return 10, nil
}

func (db *tDB) GetPostByID(id int64) (*yoradb.Post, error) {
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

func (db *tDB) CreateUser(name, email, password string) (int64, error) {
	return 1, nil
}

func (db *tDB) LoginUser(email, password string) (int64, error) {
	return 1, nil
}

func (db *tDB) GetUserByID(id int64) (*yoradb.User, error) {
	for _, v := range sessions {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

func (db *tDB) GetUserByEmail(email string) (*yoradb.User, error) {
	for _, v := range sessions {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

func (db *tDB) CreateSession(s *yoradb.Session) error {
	return nil
}

func (db *tDB) GetSessionByID(id string) (*yoradb.Session, error) {
	return &yoradb.Session{ID: id}, nil
}

func (db *tDB) UpdateSession(s *yoradb.Session) error {
	return nil
}

func (db *tDB) DeleteSession(s *yoradb.Session) error {
	return nil
}

func (db *tDB) LogoutFromSession(s *yoradb.Session) error {
	return nil
}

func (db *tDB) CreateRestorePasswordID(email, token string) (string, error) {
	return token, nil
}

func (db *tDB) GetEmailByRestoreToken(token string) (string, error) {
	return "a1@b.c", nil
}

func (db *tDB) UpdatePasswordByRestoreToken(token, email, password string) error {
	return nil
}

func (db *tDB) Close() error {
	return nil
}
