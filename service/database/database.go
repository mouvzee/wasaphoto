/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	//USER

	//Create user in the database
	CreateUser(u User) (User, error)
	//check if the userID exist and return the user, error otherwise
	CheckID(userID int) (User, error)
	//search an user by username
	GetUserByUsername(username string) (User, error)
	//get user by ID
	GetUsernamebyID(UserID int) (User, error)
	//set the username of the user
	SetUsername(username string) error
	//check if the username is already token
	CheckUsername(Username string) (bool, error)
	//change the username with a new one
	ChangeUsername(userID int, newUsername string) error
	//delete the user and his posts from the database
	DeleteUser(userID int) error

	//PHOTO

	//Create post in the database
	CreatePhoto(y Photo, ImageData []byte) (Photo, error)
	//Get the ID of the last photo uploaded
	GetLastPhotoID(userID int) (int, error)
	//Delete a photo in the user profile
	Delete_Photo(UserID int, PhotoID int) error

	//COMMENT

	//create a comment in the database
	CreateComment(commentID, creatorID, PhotoID int) error
	//delete a comment in the database
	DeleteComment(commentID, creatorID, photoID int) error

	//LIKE

	//create a like in the database
	CreateLike(userID int, creatorID int, postID int) error
	//delete a like in the database
	DeleteLike(userID, creatorID, postID int) error

	//FOLLOW

	//create a follow in the database
	CreateFollow(followerID, followedID int) error
	//get the followers of a user and return them
	GetFollowers(followedID int, offset int, limit int) ([]User, error)
	//get the following of a user and return them
	GetFollowings(followerID int, offset int, limit int) ([]User, error)
	//delete a follow in the database
	DeleteFollow(followerID, followedID int) error

	//BAN

	//create a ban in the database
	CreateBan(bannerID, bannedID int) error
	//get all the banned users from te current user
	GetBan(bannerID, offset, limit int) ([]User, error)
	//delete te ban between two users in the database
	DeleteBan(bannerID, bannedID int) error

	//PROFILE
	//visualize the profile of the user

	Ping() error
}

type appdbimpl struct {
	c   *sql.DB
	ctx context.Context
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	var err error
	// Check if tables exists. If not, the database is empty, and we need to create the structure and the tables

	_, err = db.Exec(user_table)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(post_table)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(like_table)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(comment_table)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(follow_table)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	_, err = db.Exec(ban_table)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		c:   db,
		ctx: context.Background(),
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
