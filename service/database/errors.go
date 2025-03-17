package database

import (
	"errors"
)

// ErrNoRowsAffected is returned when no rows are affected by an operation
var ErrNoRowsAffected = errors.New("no rows affected")
