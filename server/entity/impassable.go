package entity

import "github.com/google/uuid"

// ImpassableEntity represents a non-renderable entity that blocks movement.
type ImpassableEntity struct {
    Entity
}

// NewImpassableEntity creates a new ImpassableEntity at the given position.
func NewImpassableEntity(x, y int) *ImpassableEntity {
    return &ImpassableEntity{
        Entity: Entity{
			ID: uuid.New(),
			Type: "impassable",
			X: x,
			Y: y,
		},
    }
}

// GetSnapshotData returns data to include in snapshots (minimal since it's not rendered).
func (e *ImpassableEntity) GetSnapshotData() interface{} {
    return struct {
        Name string `json:"name"`
    }{
        Name: "impassable",
    }
}