package repository

import (
	"context"
	"notion/internal/models/blocks"
	"notion/internal/models/user"
	"notion/internal/models/workspace"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	usersTable     = "users"
	usersWorkspace = "workspaces"
	usersBlocks    = "blocks"
)

type Authorization interface {
	CreateUser(ctx context.Context, user user.SignUpRequest) (uuid.UUID, error)
	GetUser(ctx context.Context, user user.SignInRequest) (user.AuthUser, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type WorkspaceRepository interface {
	Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error)
	GetWorkspaces(ctx context.Context, id uuid.UUID) ([]workspace.Workspace, error)
	GetByID(ctx context.Context, id uuid.UUID) (workspace.Workspace, error)
	GetByWorkspaceID(ctx context.Context, id uuid.UUID) ([]blocks.Block, error)
	Update(ctx context.Context, name workspace.CreateWorkspaceRequest) (workspace.Workspace, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type workspaceRepository struct {
	db     *pgxpool.Pool
	client *redis.Client
}

func NewWorkspaceRepository(db *pgxpool.Pool, client *redis.Client) WorkspaceRepository {
	return &workspaceRepository{
		db:     db,
		client: client,
	}
}

type BlockRepository interface {
}

type blockRepository struct {
	db     *pgxpool.Pool
	client *redis.Client
}

func NewBlockRepository(db *pgxpool.Pool, client *redis.Client) BlockRepository {
	return &blockRepository{
		db:     db,
		client: client,
	}
}

type Repository struct {
	Authorization
	BlockRepository
	WorkspaceRepository
}

func NewRepository(db *pgxpool.Pool, client *redis.Client) *Repository {
	return &Repository{
		Authorization:       NewAuthPostgres(db, client),
		BlockRepository:     NewBlockRepository(db, client),
		WorkspaceRepository: NewWorkspaceRepository(db, client),
	}
}
