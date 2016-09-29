package yoradb

import (
	"html/template"
	"log"
	"time"
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

// DBInsertPost create new post in db
func (db *mysqlDB) DBCreatePost(post *Post, userID int) (int, error) {
	res, err := db.Conn.Exec(
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

// DBGetPostByID find post byspecified id
func (db *mysqlDB) DBGetPostByID(id int) (*Post, error) {
	row := db.Conn.QueryRow(
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
func (db *mysqlDB) DBUpdatePost(post *Post) error {
	_, err := db.Conn.Exec(
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

// DBGetPosts returns fixed number of posts
func (db *mysqlDB) DBGetPosts(num, offset int) ([]Post, error) {

	posts := make([]Post, 0, 10)

	rows, err := db.Conn.Query(
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
