package database

import (
	"database/sql"
	"errors"
)

var GetLastPhotoQUERY = `SELECT MAX(PhotoID) FROM Photo WHERE UserID= ?;`

func (db *appdbimpl) GetLastPhotoID(userID int) (int, error) {
	var _photoID = sql.NullInt64{Int64: 0, Valid: false}
	var photoID int = 0
	res, err := db.c.Query(GetLastPhotoQUERY, userID)
	if err != nil {
		return 0, err
	}
	defer func() { err = res.Close() }()

	for res.Next() {
		if err := res.Err(); err != nil {
			return 0, err
		}

		err = res.Scan(&_photoID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}

		if !_photoID.Valid {
			photoID = 0
		} else {
			photoID = int(_photoID.Int64)
		}
	}

	return photoID, err
}