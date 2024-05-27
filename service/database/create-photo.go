package database

import (
	"database/sql"
	"os"
)

var CreatePhotoQUERY = `INSERT Post(postID, userID, caption) WITH VALUES (?,?,?)`

func (db *appdbimpl) CreatePhoto(y Photo, ImageData []byte) (Photo, error) {
	_photoID, err := db.GetLastPhotoID(y.User.UserID)
	if err != nil {
		return y, err
	}

	transition, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return y, err
	}

	defer func() {
		if err != nil {
			err = transition.Rollback()
		}
		err = transition.Commit()
	}()

	y.PhotoID = _photoID + 1
	userID := y.User.UserID
	URL := GetPhotoPath(y.PhotoID, userID)

	//saving the image
	err = os.WriteFile(URL, ImageData, 0666)
	if err != nil {
		return y, err
	}

	//crop the image?????

	_, err = db.c.Exec(CreatePhotoQUERY, y.PhotoID, y.User.UserID, y.Caption)
	if err != nil {
		return y, err
	}

	new, err := db.ViewPosts(y.User.UserID, 0, 1)
	if err != nil {
		return y, err
	}

	return new[0], err
}
