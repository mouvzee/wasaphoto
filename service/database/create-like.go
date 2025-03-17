package database

var createLikeQUERY = `INSERT INTO Like (userID, creatorID, PhotoID) VALUES (?, ?, ?)`

func (db *appdbimpl) CreateLike(userID int, creatorID int, PhotoID int) error {
	_, err := db.c.Exec(createLikeQUERY, userID, creatorID, PhotoID)

	return err
}
