package database

import (
	"time"
)

type User struct {
	UserID   int
	Username string
}

type Photo struct {
	PhotoID    int
	User       User
	URL        string
	ImageData  []byte
	Caption    string
	Nlike      int
	Ncomment   int
	Liked      bool
	Created_At time.Time
}

type Comment struct {
	CommentID  int
	User       User
	Lyric      string
	Created_At time.Time
}

type Profile struct {
	User       User
	Name       string
	Follower   int
	Following  int
	IsFollowed bool
}
