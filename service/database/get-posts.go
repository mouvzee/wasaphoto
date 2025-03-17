package database

import (
	"github.com/mouvzee/wasaphoto/service/api/methods"
)

var query_GETPOSTS = `SELECT PhotoID, userID, caption, timestamp FROM Post WHERE userID=? ORDER BY timestamp DESC LIMIT ?, ?`
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

		// Get post data
		err = rows.Scan(&post.PhotoID, &user.UserID, &post.Caption, &post.Created_At)
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
		if like == 1 {
			post.Liked = true
		} else {
			post.Liked = false
		}

		// Get owner data
		user, err = db.GetUsernamebyID(user.UserID)
		if err != nil {
			return nil, err
		}

		// Set user data
		post.User = user

		// Set image path
		post.ImageData = []byte(methods.GetPostPhotoPath(user.UserID, post.PhotoID))

		posts = append(posts, post)
	}
	return posts, err
}
