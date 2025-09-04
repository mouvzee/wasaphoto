package database

import (
	"database/sql"
	"errors"
	"time"
)

var query_CREATECOMMENT = `INSERT INTO Comment (commentID, userID, textComment, PhotoID, created_at) VALUES (?, ?, ?, ?, datetime('now'));`

func (db *appdbimpl) CreateComment(userID int, photoID int, commentText string) (Comment, error) {
	var comment Comment

	// Get the last commentID
	var lastCommentID int
	lastCommentID, err := db.GetLastCommentID(photoID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return comment, err
	}

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

	_, err = tx.Exec(query_CREATECOMMENT, lastCommentID+1, userID, commentText, photoID)
	if err != nil {
		return comment, err
	}

	user, err := db.GetUsernamebyID(userID)
	if err != nil {
		return comment, err
	}

	comment = Comment{
		CommentID:  lastCommentID + 1,
		User:       user,
		PhotoID:    photoID,
		Lyric:      commentText,
		Created_At: time.Now(),
	}

	return comment, err
}
