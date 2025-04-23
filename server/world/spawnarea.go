package world

import (
	"math/rand"
)

// SpawnArea represents an area where entities can spawn
type SpawnArea struct {
    Positions      [][2]int // Valid spawn positions [x,y]
    EntitySpawnType string   // Type of entity to spawn
}

// NewSpawnArea creates a new SpawnArea
func NewSpawnArea(entitySpawnType string) *SpawnArea {
    return &SpawnArea{
        Positions:      make([][2]int, 0),
        EntitySpawnType: entitySpawnType,
    }
}

// AddPosition adds a valid spawn position to the SpawnArea
func (sa *SpawnArea) AddPosition(x, y int) {
    sa.Positions = append(sa.Positions, [2]int{x, y})
}

// GetRandomPosition returns a random spawn position from the area
func (sa *SpawnArea) GetRandomPosition() (int, int, bool) {
    if len(sa.Positions) == 0 {
        return 0, 0, false
    }
    idx := rand.Intn(len(sa.Positions))
    pos := sa.Positions[idx]
    return pos[0], pos[1], true
}