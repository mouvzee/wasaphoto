package api

import (
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

type Photo struct {
	photoID    int
	User       User
	URL        string
	caption    string
	nlike      int
	ncomment   int
	liked      bool
	created_at time.Time
}

func (y *Photo) AddingUser(user User) Photo {
	return Photo{
		photoID:    y.photoID,
		User:       User{userID: user.userID, username: user.username},
		URL:        y.URL,
		caption:    y.caption,
		nlike:      y.nlike,
		ncomment:   y.ncomment,
		liked:      y.liked,
		created_at: y.created_at,
	}
}

func (y *Photo) SavingPhoto() database.Photo {
	return database.Photo{
		photoID:    y.photoID,
		User:       database.User{user: y.User.userID, username: y.User.username},
		URL:        y.URL,
		caption:    y.caption,
		nlike:      y.nlike,
		ncomment:   y.ncomment,
		liked:      y.liked,
		created_at: y.created_at,
	}
}
