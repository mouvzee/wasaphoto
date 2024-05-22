package api

import (
	"regexp"
	"time"

	"github.com/mouvzee/wasaphoto/service/database"
)

type Photo struct {
	PhotoID int
	User    User
	URL     string
	//ImageData  []byte
	Caption    string
	Nlike      int
	Ncomment   int
	Liked      bool
	Created_At time.Time
}

func (y *Photo) AddingUser(dbUser database.User) Photo {
	return Photo{
		PhotoID: y.PhotoID,
		User:    User{UserID: dbUser.UserID, Username: dbUser.Username},
		URL:     y.URL,
		//ImageData:  y.ImageData,
		Caption:    y.Caption,
		Nlike:      y.Nlike,
		Ncomment:   y.Ncomment,
		Liked:      y.Liked,
		Created_At: y.Created_At,
	}
}

func (y *Photo) SavingPhoto() database.Photo {
	return database.Photo{
		PhotoID: y.PhotoID,
		User:    database.User{UserID: y.User.UserID, Username: y.User.Username},
		URL:     y.URL,
		//ImageData:  y.ImageData,
		Caption:    y.Caption,
		Nlike:      y.Nlike,
		Ncomment:   y.Ncomment,
		Liked:      y.Liked,
		Created_At: y.Created_At,
	}
}

func (y *Photo) TakingPhoto(dbPhoto database.Photo) error {
	imm := database.GetPath(dbPhoto.User.UserID, dbPhoto.PhotoID)

	var u User
	err := u.TakeUser(dbPhoto.User)
	if err != nil {
		return err
	}

	y.PhotoID = dbPhoto.PhotoID
	y.User = u
	y.URL = imm
	y.Caption = dbPhoto.Caption
	y.Nlike = dbPhoto.Nlike
	y.Ncomment = dbPhoto.Ncomment
	y.Liked = dbPhoto.Liked
	y.Created_At = dbPhoto.Created_At
	return nil
}

func (y *Photo) isValid() bool {
	caption := y.Caption
	validCaption := regexp.MustCompile(`^[^\/\\]{0,64}$`)
	return validCaption.MatchString(caption)
}
