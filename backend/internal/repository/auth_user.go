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
	"github.com/redis/go-redis/v9"
)

type AuthPostgres struct {
	db     *pgxpool.Pool
	client *redis.Client
}

var ErrUserAlreadyExists = errors.New("user already exists")

func NewAuthPostgres(db *pgxpool.Pool, client *redis.Client) *AuthPostgres {
	return &AuthPostgres{
		db:     db,
		client: client,
	}
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

func (a *AuthPostgres) UpdateUserName(ctx context.Context, id uuid.UUID, NewName string) error {
	query := `UPDATE users SET name = $1 WHERE id = $2`
	res, err := a.db.Exec(ctx, query, NewName, id)

	if err != nil {
		return fmt.Errorf("failed to update user name: %w", err)
	}
	if res.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	key := "user:" + id.String()
	if err := a.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("user name updated in DB, but failed to invalidate cache: %w", err)
	}

	return nil
}
