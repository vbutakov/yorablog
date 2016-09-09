package main

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Post data structure
type Post struct {
	ID          int
	Title       template.HTML
	Description template.HTML
	ImageURL    string
	Annotation  template.HTML
	Text        template.HTML
	Author      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User data structure
type User struct {
	ID               int
	Name             string
	Email            string
	Password         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CreatePostPermit bool
	EditPostPermit   bool
}

var (
	// DBConnection - connection to the database
	DBConnection *sql.DB

	// ErrLoginFailed - error for login process
	ErrLoginFailed = errors.New("User login failed")
)

// InitDB initialize database connection
func InitDB() (err error) {
	DBConnection, err = sql.Open("mysql", BaseDSN)
	if err != nil {
		return err
	}
	DBConnection.SetMaxOpenConns(10)
	return nil
}

// DBGetPosts returns fixed number of posts
func DBGetPosts(num, offset int) ([]Post, error) {

	posts := make([]Post, 0, 10)

	rows, err := DBConnection.Query(
		`SELECT p.id, p.Title, p.Description, p.ImageURL, p.Annotation, p.PostText,
         u.Name AS AuthorName,
         p.CreatedAt, p.UpdatedAt
         FROM Posts p INNER JOIN Users u ON p.Author = u.id
         ORDER BY id desc
         LIMIT ? OFFSET ?;`,
		num, offset)
	if err != nil {
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var ID int
		var Title string
		var Description string
		var ImageURL string
		var Annotation string
		var Text string
		var Author string
		var CreatedAt time.Time
		var UpdatedAt time.Time

		err = rows.Scan(&ID, &Title, &Description, &ImageURL, &Annotation, &Text, &Author, &CreatedAt, &UpdatedAt)
		if err != nil {
			log.Printf("Error in row scan inside DBGetPosts: %v\n", err)
		}
		post := &Post{ID: ID,
			Title:       template.HTML(Title),
			Description: template.HTML(Description),
			ImageURL:    ImageURL,
			Annotation:  template.HTML(Annotation),
			Text:        template.HTML(Text),
			Author:      Author,
			CreatedAt:   CreatedAt,
			UpdatedAt:   UpdatedAt}

		posts = append(posts, *post)
	}

	return posts, nil
}

// DBGetPostByID find post byspecified id
func DBGetPostByID(id int) (*Post, error) {
	row := DBConnection.QueryRow(
		`SELECT p.id, p.Title, p.Description, p.ImageURL, p.Annotation, p.PostText,
        u.Name AS AuthorName, p.CreatedAt, p.UpdatedAt
        FROM Posts p INNER JOIN Users u ON p.Author = u.id
        WHERE p.id = ?;`,
		id)

	var ID int
	var Title string
	var Description string
	var ImageURL string
	var Annotation string
	var Text string
	var Author string
	var CreatedAt time.Time
	var UpdatedAt time.Time

	err := row.Scan(&ID, &Title, &Description, &ImageURL, &Annotation, &Text, &Author, &CreatedAt, &UpdatedAt)
	post := &Post{ID: ID,
		Title:       template.HTML(Title),
		Description: template.HTML(Description),
		ImageURL:    ImageURL,
		Annotation:  template.HTML(Annotation),
		Text:        template.HTML(Text),
		Author:      Author,
		CreatedAt:   CreatedAt,
		UpdatedAt:   UpdatedAt}

	return post, err
}

// DBUpdatePost updates post in the database
func DBUpdatePost(post *Post) error {
	_, err := DBConnection.Exec(
		`UPDATE Posts
		SET
			Title = ?,
			Description = ?,
			ImageURL = ?,
			Annotation = ?,
			PostText = ?
		WHERE id = ?;`,
		string(post.Title),
		string(post.Description),
		post.ImageURL,
		string(post.Annotation),
		string(post.Text),
		post.ID)

	return err
}

// DBInsertPost create new post in db
func DBInsertPost(post *Post, userID int) (int, error) {
	res, err := DBConnection.Exec(
		`INSERT INTO Posts
			(Title, Description, ImageURL, Annotation, PostText, Author)
		VALUES(?, ?, ?, ?, ?, ?);`,
		string(post.Title),
		string(post.Description),
		post.ImageURL,
		string(post.Annotation),
		string(post.Text),
		userID)
	if err != nil {
		return 0, err
	}

	var postID int64
	postID, err = res.LastInsertId()

	return int(postID), err
}

// DBSessionValid check that session exist and not expire
func DBSessionValid(sessionID string) bool {
	var s string
	row := DBConnection.QueryRow(
		`SELECT id FROM Sessions WHERE id = ?;`,
		sessionID)

	err := row.Scan(&s)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

// DBInsertNewSession inserts new session into db
func DBInsertNewSession(sessionID string, expires time.Time) error {
	_, err := DBConnection.Exec(
		`INSERT INTO Sessions
			(id, Expires)
		VALUES (?, ?);`,
		sessionID, expires)
	return err
}

// DBUserIsLogedIn checks if user is loged in
func DBUserIsLogedIn(sessionID string) bool {
	var u sql.NullString
	row := DBConnection.QueryRow(
		`SELECT UserId FROM Sessions WHERE id = ?;`,
		sessionID)
	err := row.Scan(&u)
	if err != nil {
		log.Printf("Error in DBUserIsLogedIn: %v\n", err)
		return false
	}

	return u.Valid
}

// DBCreateUser add new user to Users table
func DBCreateUser(name, email, password string) (int, error) {
	passwordHash := getPasswordHash(email, password)
	res, err := DBConnection.Exec(
		`INSERT INTO Users
			(Name, Email, Password)
		VALUES(?,?,?);`,
		name, email, passwordHash)

	if err != nil {
		return 0, err
	}

	var id int64
	id, err = res.LastInsertId()
	return int(id), err
}

func getPasswordHash(email, password string) string {
	data := []byte(email + ":" + password)
	ph := fmt.Sprintf("%x", sha1.Sum(data))
	return ph
}

// DBUpdateSessionWithUserID link userID with sessionID
func DBUpdateSessionWithUserID(sessionID string, userID int) error {
	_, err := DBConnection.Exec(
		`UPDATE Sessions
		SET UserId = ?
		WHERE id = ?;`,
		userID, sessionID)
	return err
}

// DBLoginUser check user password in db and return userID
func DBLoginUser(email, password string) (int, error) {
	passwordHash := getPasswordHash(email, password)
	row := DBConnection.QueryRow(
		`SELECT id FROM Users WHERE Email = ? AND Password = ?;`,
		email, passwordHash)

	var userID int
	err := row.Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, ErrLoginFailed
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}

// DBGetUserBySessionID return user data for session
func DBGetUserBySessionID(sessionID string) (*User, error) {
	row := DBConnection.QueryRow(
		`SELECT
			u.id,
			u.Name,
			u.Email,
			u.Password,
			u.CreatedAt,
			u.UpdatedAt,
			u.CreatePostPermit,
			u.EditPostPermit
		FROM Users u INNER JOIN Sessions s ON u.id = s.userID
		WHERE s.id = ?;`,
		sessionID)
	var ID int
	var Name string
	var Email string
	var Password string
	var CreatedAt time.Time
	var UpdatedAt time.Time
	var CreatePostPermit bool
	var EditPostPermit bool

	err := row.Scan(&ID, &Name, &Email, &Password, &CreatedAt, &UpdatedAt,
		&CreatePostPermit, &EditPostPermit)
	user := &User{ID: ID,
		Name:             Name,
		Email:            Email,
		Password:         Password,
		CreatedAt:        CreatedAt,
		UpdatedAt:        UpdatedAt,
		CreatePostPermit: CreatePostPermit,
		EditPostPermit:   EditPostPermit}
	return user, err
}

// DBLogoutUserFromSession clears userID for session
func DBLogoutUserFromSession(sessionID string) error {
	_, err := DBConnection.Exec(
		`UPDATE Sessions
		SET UserId = NULL
		WHERE id = ?;`,
		sessionID)
	return err
}

// DBEmailExist check if user withspecified email exist in db
func DBEmailExist(email string) bool {
	row := DBConnection.QueryRow(
		`SELECT u.Email FROM Users u WHERE u.Email = ?;`,
		email)

	var s string
	err := row.Scan(&s)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	} else {
		return true
	}
}

// DBCreateRestorePasswordID create restore token and return it
func DBCreateRestorePasswordID(email string) (string, error) {
	token := CreateSessionID()
	_, err := DBConnection.Exec(
		`INSERT INTO RestorePasswords (id, Email)
        VALUES(?,?);`,
		token, email)
	return token, err
}

// DBGetEmailByRestoreToken return email for specified restore password token
func DBGetEmailByRestoreToken(token string) (string, error) {
	res := DBConnection.QueryRow(
		`SELECT Email FROM RestorePasswords WHERE id = ?;`,
		token)

	var email string
	err := res.Scan(&email)

	return email, err
}

// DBUpdatePasswordByRestoreToken update user password by specified restore token
func DBUpdatePasswordByRestoreToken(token, email, password string) error {
	tx, err := DBConnection.Begin()
	if err != nil {
		return err
	}

	passwordHash := getPasswordHash(email, password)

	_, err = tx.Exec(
		`UPDATE Users u INNER JOIN RestorePasswords rp ON u.Email = rp.Email
    SET u.Password = ?
    WHERE rp.id = ?;`,
		passwordHash, token)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM RestorePasswords
    WHERE id = ?;`,
		token)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_ = tx.Commit()

	return err
}
