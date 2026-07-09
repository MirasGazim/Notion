package service

import (
	"context"
	"notion/internal/models/workspace"

	"github.com/google/uuid"
)

func (s *workspaceService) Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error) {
	return s.repo.Create(ctx, req)
}

func (s *workspaceService) GetWorkspaces(ctx context.Context, id uuid.UUID) ([]workspace.Workspace, error) {
	return s.repo.GetWorkspaces(ctx, id)
}
