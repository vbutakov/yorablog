package yoradb

// RestorePasswordRepository interface
type RestorePasswordRepository interface {
	CreateRestorePasswordID(email, token string) (string, error)
	GetEmailByRestoreToken(token string) (string, error)
	UpdatePasswordByRestoreToken(token, email, password string) error
}

// CreateRestorePasswordID create restore token and return it
func (db *MysqlDB) CreateRestorePasswordID(email, token string) (string, error) {
	//token := CreateSessionID()
	_, err := db.Conn.Exec(
		`INSERT INTO RestorePasswords (id, Email)
        VALUES(?,?);`,
		token, email)
	return token, err
}

// GetEmailByRestoreToken return email for specified restore password token
func (db *MysqlDB) GetEmailByRestoreToken(token string) (string, error) {
	res := db.Conn.QueryRow(
		`SELECT Email FROM RestorePasswords WHERE id = ?;`,
		token)

	var email string
	err := res.Scan(&email)

	return email, err
}

// UpdatePasswordByRestoreToken update user password by specified restore token
func (db *MysqlDB) UpdatePasswordByRestoreToken(token, email, password string) error {
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
