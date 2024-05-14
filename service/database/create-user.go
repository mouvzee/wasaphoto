package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

var adduser = `INSERT INTO User (userID, username) VALUES (?,?);`
var findID = `SELECT MAX(userID) FROM User`

func (db *appdbimpl) CreateUser(x User) (User, error) {
	var user User
	user.Username = x.Username

	//userID
	var _id = sql.NullInt64{Int64: 0, Valid: false}
	y, err := db.c.Query(findID)
	if err != nil {
		return user, err
	}

	//prepare a list of row to check the ids
	var id int
	for y.Next() {
		if y.Err() != nil {
			return user, err
		}

		err = y.Scan(&_id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return user, err
		}

		//check if ID is valid
		if !_id.Valid {
			id = 0
		} else {
			id = int(_id.Int64)
		}
	}

	// set the userID
	user.UserID = id + 1

	//create a folder for the user with userID in the path
	p := "./users-data/" + fmt.Sprint(user.UserID) + "/pubblications"
	if err = os.MkdirAll(p, os.ModePerm); err != nil {
		return user, err
	}

	//insert User in the database
	_, err = db.c.Exec(adduser, user.UserID, user.Username)
	if err != nil {
		return user, err
	}

	return user, nil

}
