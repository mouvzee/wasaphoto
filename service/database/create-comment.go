package database

import (
	"database/sql"
	"time"
)

var query_CREATECOMMENT = `INSERT INTO Comment (userID, lyric, PhotoID, created_at) VALUES (?, ?, ?, datetime('now'));`

func (db *appdbimpl) CreateComment(userID int, photoID int, commentText string) (Comment, error) {
	var comment Comment

	tx, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return comment, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	result, err := tx.Exec(query_CREATECOMMENT, userID, commentText, photoID)
	if err != nil {
		return comment, err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}

	user, err := db.GetUsernamebyID(userID)
	if err != nil {
		return comment, err
	}

	comment = Comment{
		CommentID:  int(commentID),
		User:       user,
		PhotoID:    photoID,
		Lyric:      commentText,
		Created_At: time.Now(),
	}

	return comment, err
}
