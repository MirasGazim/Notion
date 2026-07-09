package blocks

import (
	"time"

	"github.com/google/uuid"
)

type Block struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"`
	ParentID    *uuid.UUID             `json:"parent_id"`
	Content     map[string]interface{} `json:"content"`
	Position    float64                `json:"position"`
	WorkspaceID uuid.UUID              `json:"workspace_id"`
	CreatedBy   uuid.UUID              `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}
