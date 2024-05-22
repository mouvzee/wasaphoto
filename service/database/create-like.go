package database

var createLikeQUERY = `INSERT INTO Like (userID, creatorID, postID) VALUES (?, ?, ?)`

func (db *appdbimpl) CreateLike(userID int, creatorID int, postID int) error {
	_, err := db.c.Exec(createLikeQUERY, userID, creatorID, postID)

	return err
}
