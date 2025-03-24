package types

import (
	"fmt"
)

// SpatialHash represents the spatial partitioning system
type SpatialHash struct {
    CellSize     int
    HashTable    map[string][]IEntity
    entityToCell map[string]string // Maps entity UUID to its last known cell
}

// NewSpatialHash creates a new spatial hash with the given cell size
func NewSpatialHash(cellSize int) *SpatialHash {
    return &SpatialHash{
        CellSize:     cellSize,
        HashTable:    make(map[string][]IEntity),
        entityToCell: make(map[string]string),
    }
}

// hashCoordinates generates a unique key for a cell
func (sh *SpatialHash) hashCoordinates(x, y int) string {
    cellX := x / sh.CellSize
    cellY := y / sh.CellSize
    return fmt.Sprintf("%d,%d", cellX, cellY)
}

// Insert adds an entity to the correct cell and tracks its location
func (sh *SpatialHash) Insert(entity IEntity) {
    x, y := entity.GetPosition()
    cellKey := sh.hashCoordinates(x, y)

    sh.HashTable[cellKey] = append(sh.HashTable[cellKey], entity)
    sh.entityToCell[entity.GetUUID().String()] = cellKey
}

// Remove deletes an entity from its cell and tracking map
func (sh *SpatialHash) Remove(entity IEntity) {
    uuid := entity.GetUUID().String()
    cellKey, exists := sh.entityToCell[uuid]
    if !exists {
        return // Entity is not being tracked
    }

    entities := sh.HashTable[cellKey]

    // Filter out the entity
    for i, e := range entities {
        if e.GetUUID() == entity.GetUUID() {
            sh.HashTable[cellKey] = append(entities[:i], entities[i+1:]...)
            break
        }
    }
    delete(sh.entityToCell, uuid) // Remove from tracking map
}

// Update moves an entity if its cell changes
func (sh *SpatialHash) Update(entity IEntity) {
    x, y := entity.GetPosition()
    newCell := sh.hashCoordinates(x, y)

    uuid := entity.GetUUID().String()
    oldCell, exists := sh.entityToCell[uuid]

    // If the entity is new or changed cells, update its position
    if !exists || newCell != oldCell {
        sh.Remove(entity)     // Remove from the old cell
        sh.Insert(entity)     // Insert into the new cell
    }
}

// GetEntitiesInCell retrieves all entities in the given cell
func (sh *SpatialHash) GetEntitiesInCell(x, y int) []IEntity {
    cellKey := sh.hashCoordinates(x, y)
    return sh.HashTable[cellKey]
}

// IsOccupied checks if a specific grid position is occupied
func (sh *SpatialHash) IsOccupied(x, y int) bool {
    entities := sh.GetEntitiesInCell(x, y)
    for _, entity := range entities {
        ex, ey := entity.GetPosition()
        if ex == x && ey == y {
            return true
        }
    }
    return false
}
