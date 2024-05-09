package database

import (
	"database/sql"
	"errors"
)

var CHECK_USERNAME = "SELECT Username FROM User WHERE Username=?"

func (db *appdbimpl) CheckIfExist(Username string) (bool, error) {
	var ex string
	err := db.c.QueryRow(CHECK_USERNAME, Username).Scan(&ex)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return true, err
}
