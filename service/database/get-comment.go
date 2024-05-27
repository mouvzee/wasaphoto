package database

var getCommentsQUERY = `SELECT commentID, userID, creatorID, postID, textComment, created_at FROM Comment WHERE creatorID = ? AND postID = ? AND hidden="0" LIMIT ?,?`

func (db *appdbimpl) GetComments(creatorID int, photoID int, offset int, limit int) ([]Comment, error) {
	var c []Comment
	rows, err := db.c.Query(getCommentsQUERY, creatorID, photoID, offset, limit)
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
