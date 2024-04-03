package database

var SETUSERNAME = "INSERT INTO user (id, username) VALUES (1, ?)"

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) Set_username(username string) error {
	_, err := db.c.Exec(SETUSERNAME, username)
	return err
}
