package database

import "time"


var getStreamQUERY = `SELECT u.userID, u.username, p.PhotoID, p.image, p.caption, p.created_at 
                        FROM Post p 
                        JOIN User u ON p.userID = u.userID 
                        JOIN Follow f ON p.userID = f.followedID 
                        WHERE f.followerID = ? 
                        AND p.userID NOT IN (SELECT bannerID FROM Ban WHERE bannedID = ?) 
                        ORDER BY p.created_at DESC 
                        LIMIT ? OFFSET ?`

var query_GETLIKECOUNT_STREAM = `SELECT COUNT(*) FROM Like WHERE PhotoID = ?`
var query_GETCOMMENTCOUNT_STREAM = `SELECT COUNT(*) FROM Comment WHERE PhotoID = ?`
var query_ISLIKED_STREAM = `SELECT COUNT(*) FROM Like WHERE PhotoID = ? AND userID = ?`

func (db *appdbimpl) GetStream(userID int, offset int, limit int) ([]Photo, error) {
	var posts []Photo

	// Get the posts from the database
	res, err := db.c.Query(getStreamQUERY, userID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { err = res.Close() }()

	for res.Next() {
		if err := res.Err(); err != nil {
			return nil, err
		}
		var user User
		var post Photo
		var createdAtStr string

		err := res.Scan(&user.UserID, &user.Username, &post.PhotoID, &post.ImageData, &post.Caption, &createdAtStr)
		if err != nil {
			return nil, err
		}

		// Parse del timestamp
		post.Created_At, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}

		// Get the likes count for the post
		err = db.c.QueryRow(query_GETLIKECOUNT_STREAM, post.PhotoID).Scan(&post.Nlike)
		if err != nil {
			return nil, err
		}

		// Get the comments count for the post
		err = db.c.QueryRow(query_GETCOMMENTCOUNT_STREAM, post.PhotoID).Scan(&post.Ncomment)
		if err != nil {
			return nil, err
		}

		// Get like status
		var likeCount int
		err = db.c.QueryRow(query_ISLIKED_STREAM, post.PhotoID, userID).Scan(&likeCount)
		if err != nil {
			return nil, err
		}
		post.Liked = likeCount > 0

		// Set the user data
		post.User = user

		// Append the post to the posts slice
		posts = append(posts, post)
	}

	return posts, err
}
