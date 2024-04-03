package api

import "git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"

type User struct {
	userID   int   
	username string 
}

func (x *User) savingUser() database.User {
	return database.User{
		userID: x.userID,
		username: x.username,
	}
}

func (x *User) takingUser(dbUser database.User) error {
	x.userID = dbUser.userID
	x.username = dbUser.username
	return nil
}
