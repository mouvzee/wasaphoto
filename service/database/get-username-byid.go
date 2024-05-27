package database

var getUserByID = `SELECT userID, username FROM User WHERE username = ?;`

func (db *appdbimpl) GetUsernamebyID(UserID int) (User, error) {
	var user User
	err := db.c.QueryRow(getUserByID, UserID).Scan(&user.UserID, &user.Username)
	return user, err
}
