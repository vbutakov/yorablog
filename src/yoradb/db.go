package yoradb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Conn *sql.DB
}

// InitDB initialize database connection
func InitDB(dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return &DB{Conn: db}, nil
}
