package database

import (
	"database/sql"
	"os"
	"github.com/mouvzee/wasaphoto/"
)

var DelPhoto = "DELETE FROM Photo WHERE UserID=? AND PhotoID=?"

func (db *appdbimpl) Delete_Photo(UserID int, PhotoID int) error {
	//get the id of the user that wants to delete the photo
	user, err := db.Get_user_byID(UserID)
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
	err = os.Remove(user.GetPath(UserID, PhotoID))
	if err != nil {
		return err
	}

	return err
}
