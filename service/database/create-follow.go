package database

var createFollowQUERY = "INSERT INTO Follow(followerID, followedID) VALUES (?,?)"

func (db *appdbimpl) CreateFollow(followerID, followedID int) error {
	_, err := db.c.Exec(createFollowQUERY, followerID, followedID)
	return err
}
