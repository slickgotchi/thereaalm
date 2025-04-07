package interfaces

import (
	"github.com/google/uuid"
)

// IEntity is the core interface for all entities.
type IEntity interface {
    GetUUID() uuid.UUID
    GetType() string
    Update(dt_s float64)
    GetPosition() (int, int)
    SetPosition(x, y int)
    GetSnapshotData() interface{}
    GetZone() IZone
    SetZone(zone IZone)
    IsNextToTargetEntity(target IEntity) bool
    SetDirection(direction string)
    GetDirection() string
    SetDirectionToTargetEntity(target IEntity)
    GetDirectionToTargetEntity(target IEntity) string
    SetDirectionToTargetPosition(x, y int)
    GetDirectionToTargetPosition(x, y int) string
}


