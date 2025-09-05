package api

// This function create a new user and save it in the database

func (rt *_router) create_user(u User) (User, error) {
	dbUser, err := rt.db.CreateUser(u.SaveUser())
	if err != nil {
		return u, err
	}

	// taking the user from database to became an object
	err = u.TakeUser(dbUser)
	if err != nil {
		return u, err
	}

	return u, nil

}
