package database

var getUser = `SELECT userID, username FROM User WHERE username = ?;`

func (db *appdbimpl) GetUserByUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow(getUser, username).Scan(&user.UserID, &user.Username)
	return user, err
}
