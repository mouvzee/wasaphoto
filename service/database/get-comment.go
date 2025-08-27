package database

var getCommentsQUERY = `SELECT commentID, userID, PhotoID, textComment, created_at 
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
		var userID int
		err = rows.Scan(&comment.CommentID, &userID, &comment.User.UserID, &comment.PhotoID, &comment.Lyric, &comment.Created_At)
		if err != nil {
			return nil, err
		}
		user, err := db.CheckID(userID)
		if err != nil {
			return nil, err
		}
		comment.User = user

		c = append(c, comment)
	}

	return c, err
}
