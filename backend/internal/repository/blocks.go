package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"notion/internal/models/blocks"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *blockRepository) Create(ctx context.Context, req blocks.CreateBlockRequest, workspaceID, userID uuid.UUID) (*blocks.Block, error) {
	const op = "repository.CreateBlock"
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

func (r *blockRepository) Update(ctx context.Context, id, workspaceID, userID uuid.UUID, fields map[string]any) (*blocks.Block, error) {
	const op = "repository.UpdateBlock"
	var setParts []string
	var args []any
	i := 1
	for col, val := range fields {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	idPos := i
	workspaceIDPos := i + 1
	userIDPos := i + 2
	args = append(args, id, workspaceID, userID)

	query := fmt.Sprintf(`
		UPDATE blocks
		SET %s, updated_at = NOW()
		WHERE id = $%d
		  AND workspace_id = $%d
		  AND workspace_id IN (
			  SELECT id FROM workspaces WHERE owner_id = $%d
		  )
		RETURNING *`,
		strings.Join(setParts, ", "), idPos, workspaceIDPos, userIDPos)

	var block blocks.Block

	err := r.db.QueryRow(ctx, query, args...).Scan(
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
			return nil, ErrBlockNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &block, nil

}

func (r *blockRepository) GetByID(ctx context.Context, id, workspaceID, userID uuid.UUID) (*blocks.Block, error) {
	query := `
		SELECT id, type, parent_id, content, position, workspace_id, created_by, created_at, updated_at
		FROM blocks
		WHERE id = $1
		  AND workspace_id = $2
		  AND workspace_id IN (
			  SELECT id FROM workspaces WHERE owner_id = $3
		  )`

	var block blocks.Block
	err := r.db.QueryRow(ctx, query, id, workspaceID, userID).Scan(
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
			return nil, ErrBlockNotFound
		}
		return nil, fmt.Errorf("repository.GetByID: %w", err)
	}

	return &block, nil
}
