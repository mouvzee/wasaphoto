package api

import (
	"regexp"
	"time"

	"github.com/mouvzee/wasaphoto/service/database"
)

type Comment struct {
	CommentID  int
	photoID    int
	User       User
	Lyric      string
	Created_At time.Time
}

// check that the comment is from a certain user and save the information from database
func (c *Comment) takingComment(dbComment database.Comment) error {
	var u User
	err := u.TakeUser(dbComment.User)
	if err != nil {
		return err
	}

	c.CommentID = dbComment.CommentID
	c.photoID = dbComment.PhotoID
	c.Lyric = dbComment.Lyric
	c.User = User(dbComment.User)
	c.Created_At = dbComment.Created_At

	// 120 types for every comment
	b, err := regexp.MatchString("^/[0-9a-zA-Z]{120}/$", c.Lyric)
	if b {
		c.Lyric = dbComment.Lyric
	}

	return err
}
