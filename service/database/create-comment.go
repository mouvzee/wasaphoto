package database

var createCommentQUERY = "INSERT INTO Comment (commentID, creatorID, postID) VALUES (?,?,?)"

func (db *appdbimpl) CreateComment(commentID, creatorID, PhotoID int) error {
	_, err := db.c.Exec(createCommentQUERY, commentID, creatorID, PhotoID)
	return err
}
