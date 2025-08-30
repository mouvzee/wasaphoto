package database

import (
	"time"
)

var query_GETPOSTS = `SELECT * FROM Post WHERE userID=? ORDER BY created_at DESC`
var query_GETLIKECOUNT = `SELECT COUNT(PhotoID) FROM Like WHERE PhotoID=?`
var query_GETCOMMENTCOUNT = `SELECT COUNT(PhotoID) FROM Comment WHERE PhotoID=?`
var query_ISLIKED = `SELECT COUNT(PhotoID) FROM Like WHERE PhotoID=? AND userID=?`

func (db *appdbimpl) GetPosts(profileUserID int, requestingUserID int) ([]Photo, error) {
	// Get the posts from the database
	rows, err := db.c.Query(query_GETPOSTS, profileUserID)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	// Create the slice of posts
	var posts []Photo

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var post Photo
		var user User
		var createdAtStr string

		// Get post data
		err = rows.Scan(&post.PhotoID, &user.UserID, &post.ImageData, &post.Caption, &createdAtStr)
		if err != nil {
			return nil, err
		}

		// Convert the string in time.Time
		post.Created_At, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}

		// Get like count
		err = db.c.QueryRow(query_GETLIKECOUNT, post.PhotoID).Scan(&post.Nlike)
		if err != nil {
			return nil, err
		}

		// Get comment count
		err = db.c.QueryRow(query_GETCOMMENTCOUNT, post.PhotoID).Scan(&post.Ncomment)
		if err != nil {
			return nil, err
		}

		// Get like status - CORREZIONE: usa requestingUserID invece di profileUserID
		var like int
		err = db.c.QueryRow(query_ISLIKED, post.PhotoID, requestingUserID).Scan(&like)
		if err != nil {
			return nil, err
		}
		post.Liked = like > 0 // Meglio usare > 0 invece di == 1

		// Get owner data
		user, err = db.GetUsernamebyID(user.UserID)
		if err != nil {
			return nil, err
		}

		// Set user data
		post.User = user

		posts = append(posts, post)
	}
	return posts, err
}
