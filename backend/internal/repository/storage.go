package repository

import "errors"

var (
	ErrUserExists        = errors.New("url already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotDeleted    = errors.New("no rows deleted")
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrBlockNotUpdated   = errors.New("no rows updated")
	ErrBlockNotFound     = errors.New("block not found")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
)
