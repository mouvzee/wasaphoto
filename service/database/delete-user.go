package database

import (
	"database/sql"
	"fmt"
	"os"
)

var deleteUserQUERY = `DELETE FROM User WHERE userID=?`

func (db *appdbimpl) DeleteUser(userID int) error {

	t, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = t.Rollback()
		}
		err = t.Commit()
	}()

	_, err = t.Exec(deleteUserQUERY, userID)
	if err != nil {
		return err
	}

	err = os.RemoveAll("./user-data" + fmt.Sprint(userID) + "/pubblications")
	if err != nil {
		return err
	}

	return err
}
