package interfaces

import "time"

type IWorldManager interface {
	Now() time.Duration
	Since(startTime time.Duration) time.Duration
	SetSimulationSpeed(multiplier float64)
}