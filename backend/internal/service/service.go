package service

import (
	"context"
	"errors"
	"notion/internal/models/workspace"
	"notion/internal/repository"
)

type Service struct {
	repository.Authorization
}

type BlockService interface {
}

type blockService struct {
	repo repository.BlockRepository
}

type WorkspaceService interface {
	Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error)
}

type workspaceService struct {
	repo repository.WorkspaceRepository
}

func NewWorkspaceService(repo repository.WorkspaceRepository) WorkspaceService {
	return &workspaceService{repo: repo}
}

func NewBlockService(repo repository.BlockRepository) BlockService {
	return &blockService{repo: repo}
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
)
