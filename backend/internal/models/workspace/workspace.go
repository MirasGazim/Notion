package workspace

import (
	"notion/internal/models/blocks"
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkspaceBlocks struct {
	Workspace Workspace      `json:"workspace"`
	Blocks    []blocks.Block `json:"blocks"`
}

type CreateWorkspaceRequest struct {
	ID   uuid.UUID `json:"-"`
	Name string    `json:"name" validate:"required,min=3"`
}
