package database

import (
	"database/sql"
	"errors"
)

var CHECK_USERNAME = "SELECT username FROM User WHERE username=?"

func (db *appdbimpl) CheckIfExist(Username string) (bool, error) {
	var ex string
	err := db.c.QueryRow(CHECK_USERNAME, Username).Scan(&ex)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return ex != "", err
}
