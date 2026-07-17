package repository

import (
	"context"
	"errors"
	"fmt"
	"notion/internal/models/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

var ErrUserAlreadyExists = errors.New("user already exists")

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(ctx context.Context, u user.SignUpRequest) (uuid.UUID, error) {
	var id uuid.UUID
	const op = "repository/auth_user/CreateUser"
	query := fmt.Sprintf("INSERT INTO %s(email, username, password_hash) values($1, $2, $3) RETURNING id", usersTable)
	row := a.db.QueryRow(ctx, query, u.Email, u.Username, u.Password)
	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.UUID{}, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (a *AuthPostgres) GetUser(ctx context.Context, u user.SignInRequest) (user.AuthUser, error) {
	var userRequest user.AuthUser
	const op = "repository/auth_user/GetUser"
	query := fmt.Sprintf("SELECT id, password_hash FROM %s WHERE username=$1", usersTable)
	row := a.db.QueryRow(ctx, query, u.Username)
	if err := row.Scan(&userRequest.ID, &userRequest.Password); err != nil {
		if err == pgx.ErrNoRows {
			return user.AuthUser{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return user.AuthUser{}, fmt.Errorf("%s: %w", op, err)
	}
	return userRequest, nil
}

func (a *AuthPostgres) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tag, err := a.db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
