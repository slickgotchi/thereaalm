package entity

import (
	"thereaalm/types"

	"github.com/google/uuid"
)

// ENTITY
type Entity struct {
	ID uuid.UUID
	Type string
    X int
    Y int
    types.Zoned
}

func (e *Entity) GetUUID() uuid.UUID { return e.ID }
func (e *Entity) GetType() string          { return e.Type }
func (e *Entity) Update(dt_s float64) {}
func (e *Entity) GetPosition() (int, int) {
	return e.X, e.Y
}
func (e *Entity) SetPosition(x, y int) {
	e.X = x
	e.Y = y
}
func (e *Entity) GetCustomData() interface{} {
    return nil
}