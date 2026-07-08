package workspace

import (
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateWorkspaceRequest struct {
	ID   uuid.UUID `json:"-"`
	Name string    `json:"name" validate:"required,min=3"`
}
