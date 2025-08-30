package database

import "database/sql"

var deleteBanQUERY = `DELETE FROM Ban WHERE bannerID=? AND bannedID=?`

func (db *appdbimpl) DeleteBan(bannerID, bannedID int) error {

	// Inizia la transazione
	trans, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = trans.Rollback()
		} else {
			_ = trans.Commit()
		}
	}()

	// Cancella il ban dalla tabella Ban
	_, err = trans.Exec(deleteBanQUERY, bannerID, bannedID)
	if err != nil {
		return err
	}

	return nil
}
