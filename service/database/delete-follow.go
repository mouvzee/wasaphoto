package database

var deleteFollowQUERY = "DELETE FROM Follow WHERE followerID=? AND followedID=?"

func (db *appdbimpl) DeleteFollow(followerID, followedID int) error {
	_, err := db.c.Exec(deleteFollowQUERY, followerID, followedID)
	return err
}
