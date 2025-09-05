package database

import "errors"

var deleteCommentQUERY = `DELETE FROM Comment 
                         WHERE commentID=? 
                         AND PhotoID=? 
                         AND (
                             userID=? 
                             OR PhotoID IN (
                                 SELECT PhotoID FROM Post WHERE PhotoID=? AND userID=?
                             )
                         )`

func (db *appdbimpl) DeleteComment(commentID, photoID, userID int) error {
	result, err := db.c.Exec(deleteCommentQUERY, commentID, photoID, userID, photoID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("comment not found or permission denied")
	}

	return nil
}
