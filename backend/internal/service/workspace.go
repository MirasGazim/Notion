package service

import (
	"context"
	"notion/internal/models/workspace"
)

func (s *workspaceService) Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error) {
	return s.repo.Create(ctx, req)
}
