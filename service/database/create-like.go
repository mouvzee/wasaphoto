package database

var createLikeQUERY = `INSERT INTO Like (userID, PhotoID) VALUES (?, ?)`

func (db *appdbimpl) CreateLike(userID int, PhotoID int) error {
	_, err := db.c.Exec(createLikeQUERY, userID, PhotoID)

	return err
}
