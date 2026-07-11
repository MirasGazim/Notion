package service

import (
	"context"
	"errors"
	"notion/internal/models/blocks"
	"notion/internal/models/user"
	"notion/internal/models/workspace"
	"notion/internal/repository"

	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(ctx context.Context, user user.SignUpRequest) (uuid.UUID, error)
	GetUser(ctx context.Context, user user.SignInRequest) (user.AuthUser, error)
}

type Service struct {
	Authorization
	BlockService
	WorkspaceService
}

type BlockService interface {
}

type AuthService struct {
	repo Authorization
}

type blockService struct {
	repo repository.BlockRepository
}

type WorkspaceService interface {
	Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error)
	GetWorkspaces(ctx context.Context, id uuid.UUID) ([]workspace.Workspace, error)
	GetByID(ctx context.Context, id uuid.UUID) (workspace.Workspace, error)
	GetByWorkspaceID(ctx context.Context, id uuid.UUID) ([]blocks.Block, error)
	UpdateWs(ctx context.Context, name workspace.CreateWorkspaceRequest) (workspace.Workspace, error)
}

type workspaceService struct {
	repo repository.WorkspaceRepository
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func NewWorkspaceService(repo repository.WorkspaceRepository) WorkspaceService {
	return &workspaceService{repo: repo}
}

func NewBlockService(repo repository.BlockRepository) BlockService {
	return &blockService{repo: repo}
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization:    NewAuthService(repo.Authorization),
		BlockService:     NewBlockService(repo.BlockRepository),
		WorkspaceService: NewWorkspaceService(repo.WorkspaceRepository),
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
)
