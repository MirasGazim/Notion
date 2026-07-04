package repository

import (
	"context"
	"notion/internal/models/user"

	"github.com/google/uuid"
)

const (
	usersTable = "users"
)

type Authorization interface {
	CreateUser(ctx context.Context, user user.Request) (uuid.UUID, error)
	GetUser(ctx context.Context, username string) (user.User, error)
}
