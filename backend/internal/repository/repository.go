package repository

import (
	"context"
	"notion/internal/models/user"
	"notion/internal/models/workspace"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTable     = "users"
	usersWorkspace = "workspaces"
)

type Authorization interface {
	CreateUser(ctx context.Context, user user.SignUpRequest) (uuid.UUID, error)
	GetUser(ctx context.Context, user user.SignInRequest) (user.AuthUser, error)
}

type WorkspaceRepository interface {
	Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error)
}

type workspaceRepository struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepository(db *pgxpool.Pool) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

type BlockRepository interface {
}

type blockRepository struct {
	db *pgxpool.Pool
}

func NewBlockRepository(db *pgxpool.Pool) BlockRepository {
	return &blockRepository{db: db}
}

type Repository struct {
	Authorization
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
