package main

import (
	"database/sql"
	"html/template"
	"log"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

// Post data structure
type Post struct {
	ID          int
	Title       template.HTML
	Description template.HTML
	ImageURL    string
	Annotation  template.HTML
	Text        template.HTML
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	// DBConnection - connection to the database
	DBConnection *sql.DB
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
		`SELECT id, Title, Description, ImageURL, Annotation, PostText,
      CreatedAt, UpdatedAt
    FROM Posts
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
		var CreatedAt time.Time
		var UpdatedAt time.Time

		err = rows.Scan(&ID, &Title, &Description, &ImageURL, &Annotation, &Text, &CreatedAt, &UpdatedAt)
		if err != nil {
			log.Printf("Error in row scan inside DBGetPosts: %v\n", err)
		}
		post := &Post{ID: ID,
			Title:       template.HTML(Title),
			Description: template.HTML(Description),
			ImageURL:    ImageURL,
			Annotation:  template.HTML(Annotation),
			Text:        template.HTML(Text),
			CreatedAt:   CreatedAt,
			UpdatedAt:   UpdatedAt}

		posts = append(posts, *post)
	}

	return posts, nil
}

// DBGetPostByID find post byspecified id
func DBGetPostByID(id int) (*Post, error) {
	row := DBConnection.QueryRow(
		`SELECT id, Title, Description, ImageURL, Annotation, PostText,
      CreatedAt, UpdatedAt
    FROM Posts
    WHERE id = ?;`,
		id)
	var ID int
	var Title string
	var Description string
	var ImageURL string
	var Annotation string
	var Text string
	var CreatedAt time.Time
	var UpdatedAt time.Time

	err := row.Scan(&ID, &Title, &Description, &ImageURL, &Annotation, &Text, &CreatedAt, &UpdatedAt)
	post := &Post{ID,
		template.HTML(Title),
		template.HTML(Description),
		ImageURL,
		template.HTML(Annotation),
		template.HTML(Text),
		CreatedAt,
		UpdatedAt}

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
func DBInsertPost(post *Post) (int, error) {
	res, err := DBConnection.Exec(
		`INSERT INTO Posts
			(Title, Description, ImageURL, Annotation, PostText)
		VALUES(?, ?, ?, ?, ?);`,
		string(post.Title),
		string(post.Description),
		post.ImageURL,
		string(post.Annotation),
		string(post.Text))
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
