package database

var query_UPDATECAPTION = `UPDATE Post SET caption = ? WHERE PhotoID = ? AND userID = ?`

func (db *appdbimpl) UpdateCaption(userID int, PhotoID int, newCaption string) error {
	_, err := db.c.Exec(query_UPDATECAPTION, newCaption, PhotoID, userID)
	return err
}
