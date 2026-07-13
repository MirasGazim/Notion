package repository

import (
	"context"
	"fmt"
	"notion/internal/models/blocks"
	"notion/internal/models/workspace"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *workspaceRepository) GetWorkspaces(ctx context.Context, id uuid.UUID) ([]workspace.Workspace, error) {
	ws := []workspace.Workspace{}

	query := fmt.Sprintf("SELECT id, owner_id, name, created_at FROM %s WHERE owner_id=$1", usersWorkspace)
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var w workspace.Workspace
		err := rows.Scan(&w.ID, &w.OwnerID, &w.Name, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ws, nil
}

func (r *workspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (workspace.Workspace, error) {
	var ws workspace.Workspace
	query := fmt.Sprintf("SELECT id, owner_id, name, created_at FROM %s WHERE id=$1", usersWorkspace)
	err := r.db.QueryRow(ctx, query, id).Scan(&ws.ID, &ws.OwnerID, &ws.Name, &ws.CreatedAt)
	if err != nil {
		return workspace.Workspace{}, fmt.Errorf("repository.GetByID: %w", err)
	}
	return ws, nil
}

func (r *workspaceRepository) GetByWorkspaceID(ctx context.Context, id uuid.UUID) ([]blocks.Block, error) {
	ws := []blocks.Block{}
	query := fmt.Sprintf("SELECT id, type, parent_id, content, position, workspace_id, created_by, created_at, updated_at FROM %s WHERE workspace_id=$1", usersBlocks)
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var w blocks.Block
		err := rows.Scan(&w.ID, &w.Type, &w.ParentID, &w.Content, &w.Position, &w.WorkspaceID, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}

		ws = append(ws, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ws, nil

}

func (r *workspaceRepository) Update(ctx context.Context, name workspace.CreateWorkspaceRequest) (workspace.Workspace, error) {
	var ws workspace.Workspace
	query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2 RETURNING id, owner_id, name, created_at", usersWorkspace)
	err := r.db.QueryRow(ctx, query, name.Name, name.ID).Scan(&ws.ID, &ws.OwnerID, &ws.Name, &ws.CreatedAt)
	if err != nil {
		return workspace.Workspace{}, err
	}

	return ws, nil

}

func (r *workspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM workspaces WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
