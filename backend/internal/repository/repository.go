package repository

import (
	"context"
	"notion/internal/models/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTable = "users"
)

type Authorization interface {
	CreateUser(ctx context.Context, user user.SignUpRequest) (uuid.UUID, error)
	GetUser(ctx context.Context, user user.SignInRequest) (user.AuthUser, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
