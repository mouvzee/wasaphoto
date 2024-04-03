package api

import (
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

type Photo struct {
	PhotoID    int
	User       User
	URL        string
	ImageData  []byte
	Caption    string
	Nlike      int
	Ncomment   int
	Liked      bool
	Created_At time.Time
}

func (y *Photo) AddingUser(user User) Photo {
	return Photo{
		PhotoID:    y.PhotoID,
		User:       User{UserID: user.UserID, Username: user.Username},
		URL:        y.URL,
		ImageData:  y.ImageData,
		Caption:    y.Caption,
		Nlike:      y.Nlike,
		Ncomment:   y.Ncomment,
		Liked:      y.Liked,
		Created_At: y.Created_At,
	}
}

func (y *Photo) SavingPhoto() database.Photo {
	return database.Photo{
		PhotoID:    y.PhotoID,
		User:       database.User{UserID: y.User.UserID, Username: y.User.Username},
		URL:        y.URL,
		ImageData:  y.ImageData,
		Caption:    y.Caption,
		Nlike:      y.Nlike,
		Ncomment:   y.Ncomment,
		Liked:      y.Liked,
		Created_At: y.Created_At,
	}
}
