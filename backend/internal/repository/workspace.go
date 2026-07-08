package repository

import (
	"context"
	"fmt"
	"notion/internal/models/workspace"
)

func (r *workspaceRepository) Create(ctx context.Context, req workspace.CreateWorkspaceRequest) (*workspace.Workspace, error) {
	query := fmt.Sprintf("INSERT INTO %s(owner_id, name) VALUES($1, $2) RETURNING id, owner_id, name, created_at", usersWorkspace)

	var ws workspace.Workspace
	err := r.db.QueryRow(ctx, query, req.ID, req.Name).Scan(&ws.ID, &ws.OwnerID, &ws.Name, &ws.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &ws, nil
}
