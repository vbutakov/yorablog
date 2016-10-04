package yoradb

import (
	"database/sql"
	"errors"
	"time"
)

// User data structure
type User struct {
	ID               int64
	Name             string
	Email            string
	Password         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CreatePostPermit bool
	EditPostPermit   bool
}

// UserRepository - interface for work with users
type UserRepository interface {
	CreateUser(name, email, password string) (int64, error)
	LoginUser(email, password string) (int64, error)
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
}

var (
	// ErrLoginFailed - error for login process
	ErrLoginFailed = errors.New("User login failed")
)

// CreateUser add new user to Users table
func (db *MysqlDB) CreateUser(name, email, password string) (int64, error) {
	passwordHash := GetPasswordHash(email, password)
	res, err := db.Conn.Exec(
		`INSERT INTO Users
			(Name, Email, Password)
		VALUES(?,?,?);`,
		name, email, passwordHash)

	if err != nil {
		return 0, err
	}

	var id int64
	id, err = res.LastInsertId()
	return id, err
}

// LoginUser check user password in db and return userID
func (db *MysqlDB) LoginUser(email, password string) (int64, error) {
	passwordHash := GetPasswordHash(email, password)
	row := db.Conn.QueryRow(
		`SELECT id FROM Users WHERE Email = ? AND Password = ?;`,
		email, passwordHash)

	var userID int64
	err := row.Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, ErrLoginFailed
	} else if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserByID reads user form db
func (db *MysqlDB) GetUserByID(id int64) (*User, error) {
	row := db.Conn.QueryRow(
		`SELECT
			u.Name,
			u.Email,
			u.CreatedAt,
			u.UpdatedAt,
			u.CreatePostPermit,
			u.EditPostPermit
		FROM Users u
		WHERE u.id = ?;`,
		id)

	var Name string
	var Email string
	var CreatedAt time.Time
	var UpdatedAt time.Time
	var CreatePostPermit bool
	var EditPostPermit bool

	err := row.Scan(&Name, &Email, &CreatedAt, &UpdatedAt, &CreatePostPermit, &EditPostPermit)
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:               id,
		Name:             Name,
		Email:            Email,
		CreatedAt:        CreatedAt,
		UpdatedAt:        UpdatedAt,
		CreatePostPermit: CreatePostPermit,
		EditPostPermit:   EditPostPermit}
	return u, nil
}

// GetUserByEmail reads user form db
func (db *MysqlDB) GetUserByEmail(email string) (*User, error) {
	row := db.Conn.QueryRow(
		`SELECT
			u.id,
			u.Name,
			u.Email,
			u.CreatedAt,
			u.UpdatedAt,
			u.CreatePostPermit,
			u.EditPostPermit
		FROM Users u
		WHERE u.Email = ?;`,
		email)

	var ID int64
	var Name string
	var Email string
	var CreatedAt time.Time
	var UpdatedAt time.Time
	var CreatePostPermit bool
	var EditPostPermit bool

	err := row.Scan(&ID, &Name, &Email, &CreatedAt, &UpdatedAt, &CreatePostPermit, &EditPostPermit)
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:               ID,
		Name:             Name,
		Email:            Email,
		CreatedAt:        CreatedAt,
		UpdatedAt:        UpdatedAt,
		CreatePostPermit: CreatePostPermit,
		EditPostPermit:   EditPostPermit}
	return u, nil
}
