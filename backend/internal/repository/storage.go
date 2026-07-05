package repository

import "errors"

var (
	ErrUserExists     = errors.New("url already exists")
	ErrUserNotFound   = errors.New("no rows")
	ErrUserNotDeleted = errors.New("no rows deleted")
)
