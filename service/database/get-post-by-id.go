package database

import "time"

var getPostByIDQUERY = `SELECT PhotoID, userID, image, caption, created_at FROM Post WHERE PhotoID = ?`
var getLikeCountQUERY = `SELECT COUNT(*) FROM Like WHERE PhotoID = ?`
var getCommentCountQUERY = `SELECT COUNT(*) FROM Comment WHERE PhotoID = ?`
var checkIfLikedQUERY = `SELECT COUNT(*) FROM Like WHERE PhotoID = ? AND userID = ?`

func (db *appdbimpl) GetPostByID(photoID int, requestingUserID ...int) (Photo, error) {
	var photo Photo
	var createdAtStr string

	// Recupera il post principale
	row := db.c.QueryRow(getPostByIDQUERY, photoID)
	err := row.Scan(&photo.PhotoID, &photo.User.UserID, &photo.ImageData, &photo.Caption, &createdAtStr)
	if err != nil {
		return photo, err
	}

	// Converte la stringa in time.Time
	photo.Created_At, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return photo, err
	}

	// Recupera i dati dell'utente
	user, err := db.GetUsernamebyID(photo.User.UserID)
	if err != nil {
		return photo, err
	}
	photo.User = user

	// Conta i like
	err = db.c.QueryRow(getLikeCountQUERY, photoID).Scan(&photo.Nlike)
	if err != nil {
		return photo, err
	}

	// Conta i commenti
	err = db.c.QueryRow(getCommentCountQUERY, photoID).Scan(&photo.Ncomment)
	if err != nil {
		return photo, err
	}

	// Controlla se un utente specifico ha messo like (se fornito)
	if len(requestingUserID) > 0 && requestingUserID[0] > 0 {
		var likeCount int
		err = db.c.QueryRow(checkIfLikedQUERY, photoID, requestingUserID[0]).Scan(&likeCount)
		if err != nil {
			return photo, err
		}
		photo.Liked = likeCount > 0
	} else {
		photo.Liked = false
	}

	return photo, nil
}
