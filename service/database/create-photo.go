package database

import (
	"database/sql"
	"time"
)

var CreatePhotoQUERY = `INSERT INTO Post(userID, image, caption, created_at) VALUES (?,?,?,datetime('now'))`

func (db *appdbimpl) CreatePhoto(y Photo, ImageData []byte) (Photo, error) {
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

	result, err := transition.Exec(CreatePhotoQUERY, y.User.UserID, ImageData, y.Caption)
	if err != nil {
		return y, err
	}

	photoID, err := result.LastInsertId()
	if err != nil {
		return y, err
	}

	y.PhotoID = int(photoID)

	y.ImageData = ImageData
	y.Created_At = time.Now()
	y.Nlike = 0
	y.Ncomment = 0
	y.Liked = false

	return y, nil
}
