package api

import (
	"regexp"

	"github.com/mouvzee/wasaphoto/service/database"
)

type User struct {
	UserID   int		
	Username string
}

func (x *User) SaveUser() database.User {
	return database.User{
		UserID:   x.UserID,
		Username: x.Username,
	}
}

func (x *User) TakeUser(dbUser database.User) error {
	x.UserID = dbUser.UserID
	x.Username = dbUser.Username
	return nil
}

func (x *User) isValid() bool {
	username := x.Username
	b, err := regexp.MatchString("^.{1,16}$", username)
	if err != nil {
		return false
	}
	return b
}
