package database

var getLikeQUERY = "SELECT userID FROM Like WHERE PhotoID=? AND creatorID=? LIMIT ?,?"

func (db *appdbimpl) GetLike(PhotoID, creatorID, offset, limit int) ([]User, error) {
	var l []User
	lines, err := db.c.Query(getLikeQUERY, PhotoID, creatorID, offset, limit)
	if err != nil {
		return nil, err
	}

	defer func() { err = lines.Close() }()

	for lines.Next() {
		if lines.Err() != nil {
			return nil, err
		}
		var userID int
		err = lines.Scan(&userID)
		if err != nil {
			return nil, err
		}
		u, err := db.CheckID(userID)
		if err != nil {
			return nil, err
		}
		l = append(l, u)
	}

	return l, err
}
