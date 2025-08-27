package api

import "github.com/mouvzee/wasaphoto/service/database"

type Profile struct {
	User       User		
	Follower   int
	Following  int
	PostsCount int
	IsFollowed bool
}

func (p *Profile) takingProfile(dbProfile database.Profile) error {
	var u User
	err := u.TakeUser(dbProfile.User)
	if err != nil {
		return err
	}

	p.User = u
	p.Follower = dbProfile.Follower
	p.Following = dbProfile.Following
	p.PostsCount = dbProfile.PostsCount
	p.IsFollowed = dbProfile.IsFollowed

	if dbProfile.IsFollowed {
		p.IsFollowed = true
	} else {
		p.IsFollowed = false
	}

	return err
}
