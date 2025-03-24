package types

import (
	"github.com/google/uuid"
)

// type EntityUUID uuid.UUID

// IEntity is the core interface for all entities.
type IEntity interface {
    GetUUID() uuid.UUID
    GetType() string
    Update(dt_s float64)
    GetPosition() (int, int)
    SetPosition(x, y int)
    GetSnapshotData() interface{}
}

// // IMovable is for entities that can move.
// type IMovable interface {
// }



type IHarvestable interface {
    Harvest() (string, int)
    CanBeHarvested() bool
}
