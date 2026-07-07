package repository

import "errors"

var (
	ErrUserExists     = errors.New("url already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrUserNotDeleted = errors.New("no rows deleted")
)
