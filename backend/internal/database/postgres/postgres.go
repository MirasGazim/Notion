package postgres

import (
	"context"
	"fmt"
	"notion/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	repository.AuthPostgres
	db *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		AuthPostgres: *repository.NewAuthPostgres(db), // используем готовый конструктор, разыменовываем указатель
		db:           db,
	}, nil
}
