package database

var ChangeUsername = `UPDATE User SET username = ? WHERE userID = ?;`

func (db *appdbimpl) ChangeUsername(userID int, newUsername string) error {
	_, err := db.c.Exec(ChangeUsername, newUsername, userID)
	return err
}
