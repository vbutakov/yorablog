package yoradb

import (
	"crypto/sha1"
	"database/sql"
	"fmt"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// MysqlDB - object for mysql db access
type MysqlDB struct {
	Conn *sql.DB
}

// InitDB initialize database connection
func InitDB(dsn string) (*MysqlDB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return &MysqlDB{Conn: db}, nil
}

// Close db connection
func (db *MysqlDB) Close() error {
	return db.Close()
}

func getPasswordHash(email, password string) string {
	data := []byte(email + ":" + password)
	ph := fmt.Sprintf("%x", sha1.Sum(data))
	return ph
}
