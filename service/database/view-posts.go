package database

var getPostsQUERY = `SELECT postID, userID, caption, timestamp FROM Post WHERE userID=? ORDER BY timestamp DESC LIMIT ?, ?`
var getlikeQUERY = `SELECT COUNT(postID) FROM Like WHERE postID=? AND userID=?`
var getCommentQUERY = `SELECT COUNT(postID) FROM Comment WHERE postID=? AND userID=?`
var statusPhotoQUERY = `SELECT COUNT(PhotoID) FROM Like WHERE userID=? AND postID=? AND creatorID=?`

func (db *appdbimpl) ViewPosts(userID, offset, limit int) ([]Photo, error) {
	lines, err := db.c.Query(getPostsQUERY, userID, offset, limit)
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

		//information about the photo
		err = lines.Scan(&photo.PhotoID, &u.UserID, &photo.Caption, &photo.Created_At)
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
		err = db.c.QueryRow(statusPhotoQUERY, userID, photo.PhotoID, u.UserID).Scan(&check)
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

		//get path for the URL
		photo.URL = GetPhotoPath(photo.PhotoID, photo.User.UserID)

		photos = append(photos, photo)
	}
	return photos, err

}
