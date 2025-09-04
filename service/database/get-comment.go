package database

import "time"

var getCommentsQUERY = `SELECT commentID, userID, PhotoID, lyric, created_at 
							FROM Comment 
							WHERE PhotoID = ?
							ORDER BY created_at ASC`

func (db *appdbimpl) GetComments(photoID int) ([]Comment, error) {
	var c []Comment
	rows, err := db.c.Query(getCommentsQUERY, photoID)
	if err != nil {
		return nil, err
	}

	defer func() { err = rows.Close() }()

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var comment Comment
		var createdAtStr string

		err = rows.Scan(&comment.CommentID, &comment.User.UserID, &comment.PhotoID, &comment.Lyric, &createdAtStr)
		if err != nil {
			return nil, err
		}

		comment.Created_At, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}

		user, err := db.CheckID(comment.User.UserID)
		if err != nil {
			return nil, err
		}
		comment.User = user

		c = append(c, comment)
	}

	return c, err
}
