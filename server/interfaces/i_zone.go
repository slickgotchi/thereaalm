package interfaces

import (
	"github.com/google/uuid"
)

type IZone interface {
	AddEntity(e IEntity)
	RemoveEntity(e IEntity)

	GetID() int
	GetPosition() (int, int)
	
	Update(dt_s float64)
	
	IsTileOccupied(x, y int) bool
	FindNearbyEntities(x, y, radius int) []IEntity
	FindNearbyEmptyTile(x, y, radius, minGap int) (int, int, bool) 
	TryGetEmptyTileNextToTargetEntity(target IEntity) (int, int, bool)

	GetEntityByUUID(uuid uuid.UUID) IEntity 
	GetEntitiesByType(entityType string) []IEntity 
	GetEntities() []IEntity

	GetDistance(x1, y1, x2, y2 int) int

	GetWorldManager() IWorldManager
}