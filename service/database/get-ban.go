package database

var getBanQUERY = `SELECT bannedID from Ban WHERE bannerID=?`

func (db *appdbimpl) GetBan(bannerID int) ([]User, error) {
	var b []User
	lines, err := db.c.Query(getBanQUERY, bannerID)
	if err != nil {
		return nil, err
	}
	defer func() { err = lines.Close() }()

	for lines.Next() {
		if lines.Err() != nil {
			return nil, err
		}
		var bannedID int
		err = lines.Scan(&bannedID)
		if err != nil {
			return nil, err
		}

		//taking the userID of the new banned user
		u, err := db.GetUsernamebyID(bannedID)
		if err != nil {
			return nil, err
		}
		//append the new banned user in the array
		b = append(b, u)
	}

	return b, err

}
