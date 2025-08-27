package database

var getLikeQUERY = `SELECT userID FROM Like WHERE PhotoID=?`

func (db *appdbimpl) GetLike(PhotoID int) ([]User, error) {
	var l []User
	lines, err := db.c.Query(getLikeQUERY, PhotoID)
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
