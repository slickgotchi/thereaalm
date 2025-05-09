// thereaalm/world/zone.go
package world

import (
	// "log"
	"math/rand"
	"thereaalm/interfaces"
	"thereaalm/utils"
	"time"

	"github.com/google/uuid"
)

type Zone struct {
    ID       int
    Entities []interfaces.IEntity
    Width    int
    Height   int
    X        int
    Y        int
    SpatialMap *SpatialHash
    WorldManager *WorldManager // Add reference to WorldManager
    ObstacleGrid [][]bool
    ThreatLevel int
}

func NewZone(wm *WorldManager, id, width, height, x, y, cellSize int) *Zone {

    obstacleGrid := make([][]bool, height)
    for i := range obstacleGrid {
        obstacleGrid[i] = make([]bool, width)
    }

    return &Zone{
        ID:       id,
        Entities: []interfaces.IEntity{},
        Width:    width,
        Height:   height,
        X:        x,
        Y:        y,
        SpatialMap: NewSpatialHash(cellSize),
        WorldManager: wm,
        ObstacleGrid: obstacleGrid,
        ThreatLevel: 50,
    }
}

func (z *Zone) GetID() int {
    return z.ID
}

func (z *Zone) GetPosition() (int, int) {
    return z.X, z.Y
}

func (z *Zone) GetWidth() int { return z.Width }
func (z *Zone) GetHeight() int { return z.Height }

func (z *Zone) AddEntity(e interfaces.IEntity) {
    z.Entities = append(z.Entities, e)
    z.SpatialMap.Insert(e)
    e.SetZone(z)
    e.SetWorldManager(z.GetWorldManager())
}

// RemoveEntity removes an entity from the zone and updates the spatial hash
func (z *Zone) RemoveEntity(e interfaces.IEntity) {
    for i, entity := range z.Entities {
        if entity.GetUUID() == e.GetUUID() {
            // Remove from entity slice
            z.Entities = append(z.Entities[:i], z.Entities[i+1:]...)
            z.SpatialMap.Remove(e) // Remove from spatial hash
            // log.Println("Removed entity from zone")
            return
        }
    }
}

// Update processes entity movement and updates spatial hash if needed
func (z *Zone) Update(dt_s float64) {
    enemyCount := 0

    for _, e := range z.Entities {
        oldX, oldY := e.GetPosition()
        e.Update(dt_s) // Allow entity to update itself
        newX, newY := e.GetPosition()
        // log.Println(e.GetType())
        if e.GetType() == "lickquidator" {
            enemyCount++
        }

        // If entity moved, update spatial hash
        if oldX != newX || oldY != newY {
            z.SpatialMap.Update(e)
        }
    }

    // log.Println(enemyCount)
    z.ThreatLevel = int(float64(enemyCount) / 1000 * 100)
    // z.ThreatLevel = int(enemyCount/1000 * 100)
    // log.Println(z.ThreatLevel)
}

func (z *Zone) GetThreatLevel() int {
    return z.ThreatLevel
}

// checks if a world position is available
func (z *Zone) IsPositionAvailable(x, y int) bool {
    return z.SpatialMap.IsPositionAvailable(x, y) && !z.IsObstacle(x, y)
}

// FindNearbyEntities finds entities within a specified radius
func (z *Zone) FindNearbyEntities(zoneX, zoneY, radius int) []interfaces.IEntity {
    entities := []interfaces.IEntity{}

    // Iterate over neighboring cells in the spatial hash
    for dx := -radius; dx <= radius; dx++ {
        for dy := -radius; dy <= radius; dy++ {
            nearby := z.SpatialMap.GetEntitiesInCell(zoneX+dx, zoneY+dy)
            entities = append(entities, nearby...)
        }
    }
    return entities
}

// FindNearbyEmptyTile finds a random empty cell within a given radius,
// ensuring a minimum gap between the returned cell and any entities.
func (z *Zone) FindNearbyAvailablePosition(x, y, radius, minGap int) (int, int, bool) {
    var candidates []struct{ dx, dy int }

    // Populate candidate positions
    for r := 1; r <= radius; r++ {
        for dx := -r; dx <= r; dx++ {
            for dy := -r; dy <= r; dy++ {
                nx, ny := x+dx, y+dy

                // Bounds check
                if nx < z.X || ny < z.Y || nx >= z.X + z.Width || ny >= z.Y + z.Height {
                    continue
                }

                candidates = append(candidates, struct{ dx, dy int }{nx, ny})
            }
        }
    }

    if len(candidates) == 0 {
        return 0, 0, false
    }

    rand.Shuffle(len(candidates), func(i, j int) {
        candidates[i], candidates[j] = candidates[j], candidates[i]
    })

    for _, candidate := range candidates {
        nx, ny := candidate.dx, candidate.dy

        // Check the gap area around the candidate position
        isValid := true
        for gx := -minGap; gx <= minGap; gx++ {
            for gy := -minGap; gy <= minGap; gy++ {
                tx, ty := nx+gx, ny+gy

                // Skip out-of-bounds tiles in the gap check
                if tx < z.X || ty < z.Y || tx >= z.X + z.Width || ty >= z.Y + z.Height {
                    continue
                }

                if !z.IsPositionAvailable(tx, ty) {
                    isValid = false
                    break
                }
            }
            if !isValid {
                break
            }
        }

        if isValid {
            return nx, ny, true
        }
    }

    return 0, 0, false
}


// TryGetEmptyCellAdjacentToEntity attempts to find an empty adjacent cell to the given entity
// Returns (x, y, true) if an empty cell is found, (0, 0, false) if no empty cell is available
func (z *Zone) TryGetEmptyTileNextToTargetEntity(target interfaces.IEntity) (int, int, bool) {
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
        if z.IsPositionAvailable(pos[0], pos[1]) { // Rename this call later
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
func (z *Zone) GetEntityByUUID(uuid uuid.UUID) interfaces.IEntity {
    for _, entity := range z.Entities {
        if entity.GetUUID() == uuid {
            return entity
        }
    }
    return nil
}

// GetEntitiesByType retrieves all entities of a specific type
func (z *Zone) GetEntitiesByType(entityType string) []interfaces.IEntity {
    var entities []interfaces.IEntity
    for _, entity := range z.Entities {
        if entity.GetType() == entityType {
            entities = append(entities, entity)
        }
    }
    return entities
}

func (z *Zone) GetEntities() []interfaces.IEntity {
    return z.Entities;
}

// GetDistance calculates the simple distance between two points (x1, y1) and (x2, y2)
func (z *Zone) GetDistance(x1, y1, x2, y2 int) int {
    return utils.Abs(x1 - x2) + utils.Abs(y1 - y2)
}

func (z *Zone) GetWorldManager() interfaces.IWorldManager {
    return z.WorldManager
}

// SetObstacle marks a tile as impassable
func (z *Zone) AddObstacle(x, y int) {
    if z.isPositionWithinZone(x, y) {
        // obstacle grid MUST be in local zone coordinates
        z.ObstacleGrid[y-z.Y][x-z.X] = true
    }
}

func (z *Zone) RemoveObstacle(x, y int) {
    if z.isPositionWithinZone(x, y) {
        // obstacle grid MUST be in local zone coordinates
        z.ObstacleGrid[y-z.Y][x-z.X] = false
    }
}

// IsObstacle checks if a tile is impassable
func (z *Zone) IsObstacle(x, y int) bool {
    if z.isPositionWithinZone(x, y) {
        // obstacle grid MUST be in local zone coordinates
        return z.ObstacleGrid[y-z.Y][x-z.X]
    } else {
        return false
    }
}

// func (z *Zone) zoneToWorldCoords(zoneX, zoneY int) (int, int) {
//     return zoneX + z.X, zoneY + z.Y
// }

// func (z *Zone) worldToZoneCoords(worldX, worldY int) (int, int) {
//     return worldX - z.X, worldY - z.Y
// }

func (z *Zone) isPositionWithinZone(x, y int) bool {
    return x >= z.X && y >= z.Y && x < z.X + z.Width && y < z.Y + z.Height
}