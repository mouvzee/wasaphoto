package database

var deleteCommentQUERY = "DELETE FROM Comment WHERE commentID=? AND creatorID=? AND postID=?"

func (db *appdbimpl) DeleteComment(commentID, creatorID, photoID int) error {
	_, err := db.c.Exec(deleteCommentQUERY, commentID, creatorID, photoID)
	if err != nil {
		return err
	}
	return nil
}
