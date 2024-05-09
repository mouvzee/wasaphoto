package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

var checkAuth = `SELECT Token FROM Authorization WHERE Token=?`

func (db *appdbimpl) IsAuthorized(UserID int) (bool, error) {
	var b string
	err := db.c.QueryRow(checkAuth, UserID).Scan(&b)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	x, err := strconv.Atoi(b)
	if err != nil {
		fmt.Println("Token impossible to find")
	}

	if UserID == x {
		return true, err
	} else {
		return false, err
	}
}
