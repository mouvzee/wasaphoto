package database

import "time"

var getPostsQUERY = `SELECT PhotoID, userID, caption, created_at FROM Post WHERE userID=? ORDER BY created_at DESC`
var getlikeQUERY = `SELECT COUNT(PhotoID) FROM Like WHERE PhotoID=? AND userID=?`
var getCommentQUERY = `SELECT COUNT(PhotoID) FROM Comment WHERE PhotoID=? AND userID=?`
var statusPhotoQUERY = `SELECT COUNT(PhotoID) FROM Like WHERE userID=? AND PhotoID=?`

func (db *appdbimpl) ViewPosts(userID int) ([]Photo, error) {
	lines, err := db.c.Query(getPostsQUERY, userID)
	if err != nil {
		return nil, err
	}
	defer func() { err = lines.Close() }()

	var photos []Photo

	for lines.Next() {
		if lines.Err() != nil {
			return nil, err
		}
		var photo Photo
		var u User
		var createdAtStr string 

		//information about the photo
		err = lines.Scan(&photo.PhotoID, &u.UserID, &photo.Caption, &createdAtStr)
		if err != nil {
			return nil, err
		}

		// parsing timestamp
		photo.Created_At, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, err
		}

		//number of like
		err = db.c.QueryRow(getlikeQUERY, photo.PhotoID, userID).Scan(&photo.Nlike)
		if err != nil {
			return nil, err
		}

		//number of comments
		err = db.c.QueryRow(getCommentQUERY, photo.PhotoID, userID).Scan(&photo.Ncomment)
		if err != nil {
			return nil, err
		}

		//check status photo between liked or not
		var check int
		err = db.c.QueryRow(statusPhotoQUERY, userID, photo.PhotoID).Scan(&check)
		if err != nil {
			return nil, err
		}
		if check == 1 {
			photo.Liked = true
		} else {
			photo.Liked = false
		}

		//find the owner information
		user, err := db.GetUsernamebyID(userID)
		if err != nil {
			return nil, err
		}

		photo.User = user

		photos = append(photos, photo)
	}
	return photos, err

}
