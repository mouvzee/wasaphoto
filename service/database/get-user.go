package database

var getUser = `SELECT userID, username FROM user WHERE username = ?;`

func (db *appdbimpl) Get_user(username string) (User, error) {
	var user User
	err := db.c.QueryRow(getUser, username).Scan(&user.userID, &user.username)
	return user, err
}
