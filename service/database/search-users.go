package database

import (
	_ "github.com/mattn/go-sqlite3"
)

var query_GETUSERS = `SELECT userID, username FROM User WHERE username LIKE ? ORDER BY username`

func (db *appdbimpl) SearchUsers(userID int, search string) ([]User, error) {
	var users []User

	rows, err := db.c.Query(query_GETUSERS, search+"%")
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		var u User
		if err := rows.Scan(&u.UserID, &u.Username); err != nil {
			return nil, err
		}

		isBanned, err := db.IsBanned(u.UserID, userID)
		if err != nil {
			return nil, err
		}

		if !isBanned {
			users = append(users, u)
		}
	}

	return users, nil
}
