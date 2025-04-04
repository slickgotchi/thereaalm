// thereaalm/world/zone.go
package types

import (
	"math/rand"
	"thereaalm/utils"
	"time"

	"github.com/google/uuid"
)

type Zone struct {
    ID       int
    Entities []IEntity
    Width    int
    Height   int
    X        int
    Y        int
    SpatialMap *SpatialHash
}

func NewZone(id, width, height, x, y, cellSize int) *Zone {
    return &Zone{
        ID:       id,
        Entities: []IEntity{},
        Width:    width,
        Height:   height,
        X:        x,
        Y:        y,
        SpatialMap: NewSpatialHash(cellSize),
    }
}

func (z *Zone) AddEntity(e IEntity) {
    z.Entities = append(z.Entities, e)
    z.SpatialMap.Insert(e)
    e.SetZone(z)
}

// RemoveEntity removes an entity from the zone and updates the spatial hash
func (z *Zone) RemoveEntity(e IEntity) {
    for i, entity := range z.Entities {
        if entity.GetUUID() == e.GetUUID() {
            // Remove from entity slice
            z.Entities = append(z.Entities[:i], z.Entities[i+1:]...)
            z.SpatialMap.Remove(e) // Remove from spatial hash
            return
        }
    }
}

// Update processes entity movement and updates spatial hash if needed
func (z *Zone) Update(dt_s float64) {
    for _, e := range z.Entities {
        oldX, oldY := e.GetPosition()
        e.Update(dt_s) // Allow entity to update itself
        newX, newY := e.GetPosition()

        // If entity moved, update spatial hash
        if oldX != newX || oldY != newY {
            z.SpatialMap.Update(e)
        }
    }
}

// IsCellOccupied checks if a specific position in the zone is occupied
func (z *Zone) IsTileOccupied(x, y int) bool {
    return z.SpatialMap.IsTileOccupied(x, y)
}

// FindNearbyEntities finds entities within a specified radius
func (z *Zone) FindNearbyEntities(x, y, radius int) []IEntity {
    entities := []IEntity{}

    // Iterate over neighboring cells in the spatial hash
    for dx := -radius; dx <= radius; dx++ {
        for dy := -radius; dy <= radius; dy++ {
            nearby := z.SpatialMap.GetEntitiesInCell(x+dx, y+dy)
            entities = append(entities, nearby...)
        }
    }
    return entities
}

// FindNearbyEmptyCell finds a random empty cell within a given radius.
func (z *Zone) FindNearbyEmptyTile(x, y, radius int) (int, int, bool) {
	// Create a random seed based on the current time (optional, for better randomness)
	// rand.Seed(time.Now().UnixNano())

	// Create a slice of coordinates to scan, including all positions within the radius.
	var candidates []struct{ dx, dy int }

	// Populate the candidates list with all offsets within the radius.
	for r := 1; r <= radius; r++ {
		for dx := -r; dx <= r; dx++ {
			for dy := -r; dy <= r; dy++ {
				nx, ny := x+dx, y+dy

				// Ensure the position is within the bounds of the zone
				if nx < z.X || ny < z.Y || nx >= z.X + z.Width || ny >= z.Y + z.Height {
					continue
				}

				// Add the position to the list of candidate positions.
				candidates = append(candidates, struct{ dx, dy int }{nx, ny})
			}
		}
	}

	// If no candidates were found, return false
	if len(candidates) == 0 {
		return 0, 0, false
	}

	// Randomly shuffle the candidates
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	// Now, check the shuffled list for the first unoccupied cell
	for _, candidate := range candidates {
		nx, ny := candidate.dx, candidate.dy

		// Check if the cell is empty
		if !z.IsTileOccupied(nx, ny) {
			return nx, ny, true
		}
	}

	return 0, 0, false // No empty cell found
}

// TryGetEmptyCellAdjacentToEntity attempts to find an empty adjacent cell to the given entity
// Returns (x, y, true) if an empty cell is found, (0, 0, false) if no empty cell is available
func (z *Zone) TryGetEmptyTileNextToTargetEntity(target IEntity) (int, int, bool) {
    // Get entity's current position
    x, y := target.GetPosition()
    
    // Define the four adjacent positions
    adjacent := [4][2]int{
        {x - 1, y}, // left
        {x + 1, y}, // right
        {x, y - 1}, // up
        {x, y + 1}, // down
    }
    
    var emptyTiles [][2]int
    for _, pos := range adjacent {
        if !z.IsTileOccupied(pos[0], pos[1]) { // Rename this call later
            emptyTiles = append(emptyTiles, pos)
        }
    }
    
    if len(emptyTiles) == 0 {
        return 0, 0, false
    }
    
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    chosen := emptyTiles[r.Intn(len(emptyTiles))]
    return chosen[0], chosen[1], true
}

// GetEntityByUUID retrieves an entity by its UUID
func (z *Zone) GetEntityByUUID(uuid uuid.UUID) IEntity {
    for _, entity := range z.Entities {
        if entity.GetUUID() == uuid {
            return entity
        }
    }
    return nil
}

// GetEntitiesByType retrieves all entities of a specific type
func (z *Zone) GetEntitiesByType(entityType string) []IEntity {
    var entities []IEntity
    for _, entity := range z.Entities {
        if entity.GetType() == entityType {
            entities = append(entities, entity)
        }
    }
    return entities
}

// GetDistance calculates the simple distance between two points (x1, y1) and (x2, y2)
func (z *Zone) GetDistance(x1, y1, x2, y2 int) int {
    return utils.Abs(x1 - x2) + utils.Abs(y1 - y2)
}