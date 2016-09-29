package yoradb

import (
	"database/sql"
	"time"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// DB interface
type DB interface {
	DBCreatePost(post *Post, userID int) (int, error)
	DBGetPostByID(id int) (*Post, error)
	DBUpdatePost(post *Post) error
	DBGetPosts(num, offset int) ([]Post, error)

	DBGetUserBySessionID(sessionID string) (*User, error)
	DBSessionValid(sessionID string) bool
	DBInsertNewSession(sessionID string, expires time.Time) error
	DBCreateUser(name, email, password string) (int, error)
	DBUpdateSessionWithUserID(sessionID string, userID int) error
	DBLoginUser(email, password string) (int, error)
	DBLogoutUserFromSession(sessionID string) error
	DBEmailExist(email string) bool
	DBCreateRestorePasswordID(email, token string) (string, error)
	DBGetEmailByRestoreToken(token string) (string, error)
	DBUpdatePasswordByRestoreToken(token, email, password string) error

	Close() error
}

type mysqlDB struct {
	Conn *sql.DB
}

// InitDB initialize database connection
func InitDB(dsn string) (DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return &mysqlDB{Conn: db}, nil
}

func (db *mysqlDB) Close() error {
	return db.Close()
}
