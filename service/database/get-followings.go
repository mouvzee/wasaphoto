package database

var getFollowingsQUERY = `SELECT userID, username FROM User WHERE userID IN (SELECT followedID FROM Follow WHERE followerID=?)`

func (db *appdbimpl) GetFollowings(followerID int) ([]User, error) {
	//taking the followings from the database
	rows, err := db.c.Query(getFollowingsQUERY, followerID)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	var f []User

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var following User

		//taking the data of the following to append it
		err = rows.Scan(&following.UserID, &following.Username)
		if err != nil {
			return nil, err
		}
		f = append(f, following)
	}

	return f, err
}
