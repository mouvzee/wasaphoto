package database

import (
	"database/sql"
	"fmt"
	"os"
	//"github.com/mouvzee/wasaphoto/service/api/methods"
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
	err = os.Remove(GetPath(UserID, PhotoID))
	if err != nil {
		return err
	}

	return err
}

func GetPath(UserID, PhotoID int) string {
	return fmt.Sprintf("./users/%d/posts/%d.jpeg", UserID, PhotoID)
}
