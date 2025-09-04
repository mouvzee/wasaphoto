package database

import (
	"database/sql"
)

var DelPhoto = "DELETE FROM Post WHERE PhotoID=?"

func (db *appdbimpl) Delete_Photo(PhotoID int) error {
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
	_, err = tx.Exec(DelPhoto, PhotoID)
	if err != nil {
		return err
	}

	return err
}
