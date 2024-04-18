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
	//Create user in the database
	Create_user(u User) (User, error)
	//search an user by username
	Get_user_byUsername(username string) (User, error)
	//get user by ID
	Get_user_byID(UserID int) (User, error)
	//set the username of the user
	Set_username(username string) error
	//Post a photo in the user profile
	Post_Photo(y Photo, ImageData []byte) error
	//Delete a photo in the user profile
	Delete_Photo(UserID int, PhotoID int) error
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

	// Check if tables exists. If not, the database is empty, and we need to create the structure and the tables
	var user_table string
	err := db.QueryRow(`SELECT user_table FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&user_table)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(user_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	var post_table string
	err = db.QueryRow(`SELECT post_table FROM sqlite_master WHERE type='table' AND name='post_table';`).Scan(&post_table)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(post_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	var like_table string
	err = db.QueryRow(`SELECT like_table FROM sqlite_master WHERE type='table' AND name='like_table';`).Scan(&like_table)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(like_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	var comment_table string
	err = db.QueryRow(`SELECT comment_table FROM sqlite_master WHERE type='table' AND name='comment_table';`).Scan(&comment_table)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(comment_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	var follow_table string
	err = db.QueryRow(`SELECT follow_table FROM sqlite_master WHERE type='table' AND name='follow_table';`).Scan(&follow_table)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(follow_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	var ban_table string
	err = db.QueryRow(`SELECT ban_table FROM sqlite_master WHERE type='table' AND name='ban_table';`).Scan(&ban_table)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(ban_table)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c:   db,
		ctx: context.Background(),
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
