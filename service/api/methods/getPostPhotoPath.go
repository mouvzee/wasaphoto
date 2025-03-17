package methods

import (
	"fmt"
)

func GetPostPhotoPath(ownerID int, PhotoID int) string {
	return fmt.Sprintf("./storage/%d/posts/%d.jpeg", ownerID, PhotoID)
}
