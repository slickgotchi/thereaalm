package interfaces

import (
	"github.com/google/uuid"
)

type IZone interface {
	AddEntity(e IEntity)
	RemoveEntity(e IEntity)

	GetID() int
	GetPosition() (int, int)

	GetWidth() int
	GetHeight() int
	
	Update(dt_s float64)
	
	IsPositionAvailable(zoneX, zoneY int) bool
	FindNearbyEntities(zoneX, zoneY, radius int) []IEntity
	FindNearbyAvailablePosition(zoneX, zoneY, radius, minGap int) (int, int, bool) 
	TryGetEmptyTileNextToTargetEntity(target IEntity) (int, int, bool)

	GetEntityByUUID(uuid uuid.UUID) IEntity 
	GetEntitiesByType(entityType string) []IEntity 
	GetEntities() []IEntity

	GetDistance(x1, y1, x2, y2 int) int

	GetWorldManager() IWorldManager

	AddObstacle(x, y int)
	RemoveObstacle(x, y int)
	IsObstacle(x, y int) bool
}