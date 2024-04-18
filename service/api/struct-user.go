package api

import "github.com/mouvzee/wasaphoto/service/database"

type User struct {
	UserID   int
	Username string
}

func (x *User) savingUser() database.User {
	return database.User{
		UserID:   x.UserID,
		Username: x.Username,
	}
}

func (x *User) takingUser(dbUser database.User) error {
	x.UserID = dbUser.UserID
	x.Username = dbUser.Username
	return nil
}
