package methods

import (
	"fmt"
)

func GetPath(ownerID int, postID int) string {
	return fmt.Sprintf("./users/%d/posts/%d.jpeg", ownerID, postID)
}
