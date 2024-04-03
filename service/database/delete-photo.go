package database

import (
	"database/sql"
	"os"

	"/home/mouvzee/wasa/gits/wasaphoto/users"
)

var DelPhoto = "DELETE FROM Photo WHERE UserID=? AND PhotoID=?"

func (db *appdbimpl) Delete_Photo(UserID int, PhotoID int) error {
	tx, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
		err = tx.Commit()
	}()

	// Delete the post
	_, err = tx.Exec(DelPhoto, UserID, PhotoID)
	if err != nil {
		return err
	}

	// Delete file
	err = os.Remove(users.GetPostPhotoPath(UserID, PhotoID))
	if err != nil {
		return err
	}

	return err
}
