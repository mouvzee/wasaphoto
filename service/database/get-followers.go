package database

var getFollowersQUERY = `SELECT userID, username FROM User WHERE userID IN (SELECT followerID FROM Follow WHERE followedID=?)`

func (db *appdbimpl) GetFollowers(followedID int) ([]User, error) {
	//search the followers in the database to take them
	rows, err := db.c.Query(getFollowersQUERY, followedID)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	var followers []User

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var follower User

		//taking the follower data to append it
		err = rows.Scan(&follower.UserID, &follower.Username)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	return followers, err
}
