package api

import "github.com/mouvzee/wasaphoto/service/database"

type Profile struct {
	User       User
	Name       string
	Follower   int
	Following  int
	IsFollowed bool
}

func (p *Profile) takingProfile(dbProfile database.Profile) error {
	var u User
	err := u.takingUser(dbProfile.User)
	if err != nil {
		return err
	}

	p.User = u
	p.Follower = dbProfile.Follower
	p.Following = dbProfile.Following
	p.Name = dbProfile.Name
	p.IsFollowed = dbProfile.IsFollowed

	if dbProfile.IsFollowed {
		p.IsFollowed = true
	} else {
		p.IsFollowed = false
	}

	return err
}
