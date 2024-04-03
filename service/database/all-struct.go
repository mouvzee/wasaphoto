package database

import "time"

type User struct {
	userID   int
	username string
}

type Photo struct {
	photoID		int
	User		User
	URL			string
	caption		string
	nlike		int
	ncomment	int
	liked		bool
	created_at	time.Time
}