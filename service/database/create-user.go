package database

import (
	"database/sql"
	"fmt"
	"os"
)

var authUser = `INSERT INTO Authorization (User, token) VALUES (?,?)`
var adduser = `INSERT INTO User (userID, username) VALUES (?,?);`
var findID = `SELECT userID FROM User`

func (db *appdbimpl) CreateUser(x User) (User, error) {
	var user User
	user.Username = x.Username

	//userID
	var id = sql.NullInt64{Int64: 0, Valid: false}
	y, err := db.c.Query(findID)
	if err != nil {
		return user, err
	}

	//prepare a list of row to check the ids
	for y.Next() {
		if y.Err() != nil {
			return user, err
		}
		//check if ID is valid
		if !id.Valid {
			fmt.Println("Invalid ID, ID must be NOT NULL and a positive integer")
			return user, err
		}
	}

	//create a folder for the user with userID in the path
	p := "./users-data/" + fmt.Sprint(user.UserID) + "/pubblications"
	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return user, err
	}
	//create a unique token for the user
	var a Authorization
	a.User = user
	a.Token = user.UserID
	_, err = db.c.Exec(authUser, a.User, a.Token)
	if err != nil {
		return user, err
	}

	//insert User in the database
	_, err = db.c.Exec(adduser, user.UserID, user.Username)
	if err != nil {
		return user, err
	}

	return user, nil

}
