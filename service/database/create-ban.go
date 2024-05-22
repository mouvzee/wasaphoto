package database

var createBanQUERY = "INSERT INTO Ban(bannerID, bannedID) VALUES (?,?)"

func (db *appdbimpl) CreateBan(bannerID, bannedID int) error {
	_, err := db.c.Exec(createBanQUERY, bannerID, bannedID)
	return err
}
