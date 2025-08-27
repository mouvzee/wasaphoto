package database

import (
	"database/sql"
	//"os"
)

var CreatePhotoQUERY = `INSERT INTO Post(PhotoID, userID, image, caption, created_at) VALUES (?,?,?,?,datetime('now'))`

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
			_ = transition.Rollback()
		} else {
			_ = transition.Commit()
		}
	}()

	y.PhotoID = _photoID + 1

	_, err = db.c.Exec(CreatePhotoQUERY, y.PhotoID, y.User.UserID, ImageData, y.Caption)
	if err != nil {
		return y, err
	}

	new, err := db.ViewPosts(y.User.UserID)
	if err != nil {
		return y, err
	}

	return new[0], err
}
