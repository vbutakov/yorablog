package main

import (
	"database/sql"
	"log"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

// Post data structure
type Post struct {
	ID          int
	Title       string
	Description string
	ImageURL    string
	Annotation  string
	Text        string
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
		post := &Post{ID, Title, Description, ImageURL, Annotation, Text, CreatedAt, UpdatedAt}

		posts = append(posts, *post)
	}

	return posts, nil
}
