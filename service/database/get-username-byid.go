package database

var getUserByID = `SELECT userID, username FROM user WHERE username = ?;`

func (db *appdbimpl) Get_user_byID(UserID int) (User, error) {
	var user User
	err := db.c.QueryRow(getUserByID, UserID).Scan(&user.UserID, &user.Username)
	return user, err
}
