package database

import "fmt"

func GetPhotoPath(PhotoID int, UserID int) string {
	return fmt.Sprintf("./users-data/%d/pubblications/%d.jpeg", UserID, PhotoID)
}
