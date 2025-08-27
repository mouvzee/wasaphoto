package database

var deleteCommentQUERY = `DELETE FROM Comment WHERE commentID=? AND PhotoID=?`

func (db *appdbimpl) DeleteComment(commentID, photoID int) error {
	_, err := db.c.Exec(deleteCommentQUERY, commentID, photoID)
	if err != nil {
		return err
	}
	return nil
}
