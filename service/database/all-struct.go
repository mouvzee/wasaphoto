package database

import (
	"time"
)

type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
}

type Photo struct {
	PhotoID    int       `json:"photoID"`
	User       User      `json:"user"`
	URL        string    `json:"URL"`
	ImageData  []byte    `json:"imageData"`
	Caption    string    `json:"caption"`
	Nlike      int       `json:"nlike"`
	Ncomment   int       `json:"ncomment"`
	Liked      bool      `json:"liked"`
	Created_At time.Time `json:"createdAt"`
}

type Comment struct {
	CommentID  int       `json:"commentID"`
	PhotoID    int       `json:"photoID"`
	User       User      `json:"user"`
	Lyric      string    `json:"lyric"`
	Created_At time.Time `json:"createdAt"`
}

type Profile struct {
	User       User   `json:"user"`
	Name       string `json:"name"`
	Follower   int    `json:"follower"`
	Following  int    `json:"following"`
	IsFollowed bool   `json:"isFollowed"`
}

type Authorization struct {
	User  User `json:"user"`
	Token int  `json:"token"`
}
