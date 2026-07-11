package service

import (
	"context"
	"notion/internal/models/blocks"
	"notion/internal/models/workspace"

	"github.com/google/uuid"
)

func (s *workspaceService) Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error) {
	return s.repo.Create(ctx, req)
}

func (s *workspaceService) GetWorkspaces(ctx context.Context, id uuid.UUID) ([]workspace.Workspace, error) {
	return s.repo.GetWorkspaces(ctx, id)
}

func (s *workspaceService) GetByID(ctx context.Context, id uuid.UUID) (workspace.Workspace, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *workspaceService) GetByWorkspaceID(ctx context.Context, id uuid.UUID) ([]blocks.Block, error) {
	return s.repo.GetByWorkspaceID(ctx, id)
}

func (s *workspaceService) UpdateWs(ctx context.Context, name workspace.CreateWorkspaceRequest) (workspace.Workspace, error) {
	return s.repo.Update(ctx, name)
}
