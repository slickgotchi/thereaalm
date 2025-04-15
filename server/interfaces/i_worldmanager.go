package interfaces

import "time"

type IWorldManager interface {
	Now() time.Duration
	Since(startTime time.Duration) time.Duration
	SetSimulationSpeed(multiplier float64)

	AddEntity(e IEntity)
	RemoveEntity(e IEntity)

	// utility functions 
	IsPositionAvailable(x, y int) bool
	FindNearbyAvailablePosition(x, y, radius, minimumGap int) (int, int, bool)
	GetDistance(x1, y1, x2, y2 int) int
}