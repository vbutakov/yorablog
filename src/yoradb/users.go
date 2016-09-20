package yoradb

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

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
	// ErrLoginFailed - error for login process
	ErrLoginFailed = errors.New("User login failed")
)

// DBGetUserBySessionID return user data for session
func (db *mysqlDB) DBGetUserBySessionID(sessionID string) (*User, error) {
	row := db.Conn.QueryRow(
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

// DBSessionValid check that session exist and not expire
func (db *mysqlDB) DBSessionValid(sessionID string) bool {
	var s string
	row := db.Conn.QueryRow(
		`SELECT id FROM Sessions WHERE id = ?;`,
		sessionID)

	err := row.Scan(&s)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

// DBInsertNewSession inserts new session into db
func (db *mysqlDB) DBInsertNewSession(sessionID string, expires time.Time) error {
	_, err := db.Conn.Exec(
		`INSERT INTO Sessions
			(id, Expires)
		VALUES (?, ?);`,
		sessionID, expires)
	return err
}

// DBUserIsLogedIn checks if user is loged in
func (db *mysqlDB) DBUserIsLogedIn(sessionID string) bool {
	var u sql.NullString
	row := db.Conn.QueryRow(
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
func (db *mysqlDB) DBCreateUser(name, email, password string) (int, error) {
	passwordHash := getPasswordHash(email, password)
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
	return int(id), err
}

func getPasswordHash(email, password string) string {
	data := []byte(email + ":" + password)
	ph := fmt.Sprintf("%x", sha1.Sum(data))
	return ph
}

// DBUpdateSessionWithUserID link userID with sessionID
func (db *mysqlDB) DBUpdateSessionWithUserID(sessionID string, userID int) error {
	_, err := db.Conn.Exec(
		`UPDATE Sessions
		SET UserId = ?
		WHERE id = ?;`,
		userID, sessionID)
	return err
}

// DBLoginUser check user password in db and return userID
func (db *mysqlDB) DBLoginUser(email, password string) (int, error) {
	passwordHash := getPasswordHash(email, password)
	row := db.Conn.QueryRow(
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

// DBLogoutUserFromSession clears userID for session
func (db *mysqlDB) DBLogoutUserFromSession(sessionID string) error {
	_, err := db.Conn.Exec(
		`UPDATE Sessions
		SET UserId = NULL
		WHERE id = ?;`,
		sessionID)
	return err
}

// DBEmailExist check if user withspecified email exist in db
func (db *mysqlDB) DBEmailExist(email string) bool {
	row := db.Conn.QueryRow(
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
func (db *mysqlDB) DBCreateRestorePasswordID(email, token string) (string, error) {
	//token := CreateSessionID()
	_, err := db.Conn.Exec(
		`INSERT INTO RestorePasswords (id, Email)
        VALUES(?,?);`,
		token, email)
	return token, err
}

// DBGetEmailByRestoreToken return email for specified restore password token
func (db *mysqlDB) DBGetEmailByRestoreToken(token string) (string, error) {
	res := db.Conn.QueryRow(
		`SELECT Email FROM RestorePasswords WHERE id = ?;`,
		token)

	var email string
	err := res.Scan(&email)

	return email, err
}

// DBUpdatePasswordByRestoreToken update user password by specified restore token
func (db *mysqlDB) DBUpdatePasswordByRestoreToken(token, email, password string) error {
	tx, err := db.Conn.Begin()
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
