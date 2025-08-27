package database

import (
	"time"
)

// Modifica la query per selezionare anche la colonna 'image'
var query_GETPOSTS = `SELECT PhotoID, userID, image, caption, created_at FROM Post WHERE userID=? ORDER BY created_at DESC LIMIT ?, ?`
var query_GETLIKECOUNT = `SELECT COUNT(PhotoID) FROM Like WHERE PhotoID=? AND ownerID=?`
var query_GETCOMMENTCOUNT = `SELECT COUNT(PhotoID) FROM Comment WHERE PhotoID=? AND ownerID=?`
var query_ISLIKED = `SELECT COUNT(PhotoID) FROM Like WHERE PhotoID=? AND ownerID=? AND userID=?`

func (db *appdbimpl) GetPosts(userID int, profileUserID int, offset int, limit int) ([]Photo, error) {
	// Get the posts from the database
	rows, err := db.c.Query(query_GETPOSTS, profileUserID, offset, limit)
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
		err = db.c.QueryRow(query_GETLIKECOUNT, post.PhotoID, profileUserID).Scan(&post.Nlike)
		if err != nil {
			return nil, err
		}

		// Get comment count
		err = db.c.QueryRow(query_GETCOMMENTCOUNT, post.PhotoID, profileUserID).Scan(&post.Ncomment)
		if err != nil {
			return nil, err
		}

		// Get like status
		var like int
		err = db.c.QueryRow(query_ISLIKED, post.PhotoID, user.UserID, userID).Scan(&like)
		if err != nil {
			return nil, err
		}
		post.Liked = like == 1

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
