package yoradb

import "time"

// Session - user session data
type Session struct {
	ID        string
	UserID    int
	Expires   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SessionRepository - access to session storage
type SessionRepository interface {
	CreateSession(s *Session) error
	GetSessionByID(id string) (*Session, error)
	UpdateSession(s *Session) error
	DeleteSession(s *Session) error
}

// CreateSession - insert new session into db
func (db *MysqlDB) CreateSession(s *Session) error {
	_, err := db.Conn.Exec(
		`INSERT INTO Sessions
			(id, UserId, Expires)
		VALUES(?, ?, ?);`,
		s.ID,
		s.UserID,
		s.Expires)

	return err
}

// GetSessionByID - read session information from db
func (db *MysqlDB) GetSessionByID(id string) (*Session, error) {
	row := db.Conn.QueryRow(
		`SELECT s.id, s.UserId, s.Expires, s.CreatedAt, s.UpdatedAt
        FROM Sessions s
        WHERE s.id = ?;`,
		id)

	var ID string
	var UserID int
	var Expires time.Time
	var CreatedAt time.Time
	var UpdatedAt time.Time

	err := row.Scan(&ID, &UserID, &Expires, &CreatedAt, &UpdatedAt)
	s := &Session{ID: ID,
		UserID:    UserID,
		Expires:   Expires,
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt}

	return s, err
}

// UpdateSession - modify existing session data in db
func (db *MysqlDB) UpdateSession(s *Session) error {
	_, err := db.Conn.Exec(
		`UPDATE Sessions
		SET
			id = ?,
			UserId = ?,
			Expires = ?
		WHERE id = ?;`,
		s.ID,
		s.UserID,
		s.Expires,
		s.ID)

	return err
}

// DeleteSession - delete session from db
func (db *MysqlDB) DeleteSession(s *Session) error {
	_, err := db.Conn.Exec(
		`DELETE FROM Sessions
		WHERE id = ?;`,
		s.ID)

	return err
}
