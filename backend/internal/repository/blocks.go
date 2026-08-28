package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"notion/internal/models/blocks"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *blockRepository) Create(ctx context.Context, req blocks.CreateBlockRequest, workspaceID, userID uuid.UUID) (*blocks.Block, error) {
	const op = "repository.CreateBlock	"
	query := `INSERT INTO blocks (type, parent_id, content, position, workspace_id, created_by) SELECT $1, $2, $3, $4, $5, $6 WHERE EXISTS (SELECT 1 FROM workspaces WHERE id = $5 AND owner_id = $6) RETURNING id, type, parent_id, content, position, workspace_id, created_by, created_at, updated_at`
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to marshal content: %w", op, err)
	}

	var block blocks.Block
	err = r.db.QueryRow(ctx, query, req.Type, req.ParentID, contentBytes, req.Position, workspaceID, userID).Scan(
		&block.ID,
		&block.Type,
		&block.ParentID,
		&block.Content,
		&block.Position,
		&block.WorkspaceID,
		&block.CreatedBy,
		&block.CreatedAt,
		&block.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWorkspaceNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &block, nil
}
