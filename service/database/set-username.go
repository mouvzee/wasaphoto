package database

var SETUSERNAME = `INSERT INTO User (id, username) VALUES (1, ?)`

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) SetUsername(username string) error {
	_, err := db.c.Exec(SETUSERNAME, username)
	return err
}
