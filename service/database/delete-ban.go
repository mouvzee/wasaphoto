package database

import "database/sql"

var deleteBanQUERY = "DELETE FROM Ban WHERE bannerID=? AND bannedID=?"
var visualizeCommentsQUERY = "UPDATE Comment SET hidden = FALSE WHERE userID = ? AND postID = ?"

func (db *appdbimpl) DeleteBan(bannerID, bannedID int) error {

	//taking all the posts from the banner
	lines, err := db.c.Query(getPostsQUERY, bannerID)
	if err != nil {
		return err
	}

	var ids []int
	for lines.Next() {
		if lines.Err() != nil {
			return err
		}

		var photoID int
		err = lines.Scan(&photoID)
		if err != nil {
			return err
		}

		ids = append(ids, photoID)
	}

	defer func() { err = lines.Close() }()
	//finisci sta funzione x favore devi fa la transazione

	trans, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = trans.Rollback()
		}
		err = trans.Commit()
	}()

	//e poi devi mette che l'utente può rivedere di nuovo tutti i commenti

	c, err := trans.Prepare(visualizeCommentsQUERY)
	for _, photoID := range ids {
		_, err := c.Exec(bannedID, photoID)
		if err != nil {
			return err
		}
	}
	//devi cancellare la riga nella tabella ban ce implica questa relazione tra i due

	_, err = trans.Exec(deleteBanQUERY, bannerID, bannedID)
	if err != nil {
		return err
	}
	return err
}
