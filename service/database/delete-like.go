package database

var deleteLikeQUERY = "DELETE FROM Like WHERE userID=? AND creatorID=? AND PhotoID=?"

func (db *appdbimpl) DeleteLike(userID, creatorID, PhotoID int) error {
	p, err := db.c.Exec(deleteLikeQUERY, userID, creatorID, PhotoID)
	if err != nil {
		return err
	}

	i, err := p.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return nil //informa il database ce non ci sono errori
	}
	return nil
}
