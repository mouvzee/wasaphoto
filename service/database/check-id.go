package database

var CheckIDQUERY = `SELECT userID, username FROM User WHERE userID = ?;`

func (db *appdbimpl) CheckID(userID int) (User, error) {
	var user User
	err := db.c.QueryRow(CheckIDQUERY, userID).Scan(&user.UserID, &user.Username)
	return user, err
}
