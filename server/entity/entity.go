package entity

import (
	"log"
	"math"
	"thereaalm/interfaces"

	"github.com/google/uuid"
)

// ENTITY
type Entity struct {
	ID uuid.UUID
	Type string
    X int
    Y int
	CurrentZone interfaces.IZone
	Direction string
	WorldManager interfaces.IWorldManager
}

func (e *Entity) GetUUID() uuid.UUID { return e.ID }
func (e *Entity) GetType() string          { return e.Type }
func (e *Entity) Update(dt_s float64) {}
func (e *Entity) GetPosition() (int, int) {
	return e.X, e.Y
}
func (e *Entity) SetPosition(x, y int) {
	e.X = x
	e.Y = y
}
func (e *Entity) GetCustomData() interface{} {
    return nil
}

func (e *Entity) SetZone(zone interfaces.IZone) {
    e.CurrentZone = zone
}

func (e *Entity) GetZone() interfaces.IZone {
    return e.CurrentZone
}

func (e *Entity) SetWorldManager(wm interfaces.IWorldManager) {
    e.WorldManager = wm
}

func (e *Entity) GetWorldManager() interfaces.IWorldManager {
    return e.WorldManager
}

func (e *Entity) GetSnapshotData() interface {} {
	return nil
}

func (e *Entity) IsNextToTargetEntity(target interfaces.IEntity) bool {
    ax, ay := e.GetPosition()
    bx, by := target.GetPosition()

    // Check if the entities are next to each other (left, right, up, down)
    return (ax == bx && (ay == by+1 || ay == by-1)) || // Vertical check (up, down)
           (ay == by && (ax == bx+1 || ax == bx-1))   // Horizontal check (left, right)
}

func (e *Entity) SetDirectionToTargetEntity(target interfaces.IEntity) {
	e.SetDirection(e.GetDirectionToTargetEntity(target))
}

func (e *Entity) GetDirectionToTargetEntity(target interfaces.IEntity) string {
	ax, ay := e.GetPosition()
	bx, by := target.GetPosition()

	dx := bx - ax
	dy := by - ay

	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		if dx < 0 {
			return "left"
		} else {
			return "right"
		}
	} else {
		if dy < 0 {
			return "up"
		} else {
			return "down"
		}
	}
}

func (e *Entity) SetDirection(direction string) {
	if direction == "left" || direction == "right" || 
	direction == "up" || direction == "down" {
		e.Direction = direction	
	} else {
		log.Printf("ERROR: Invalid direction '%s' passed to setDirection()", direction)
	}
}

func (e *Entity) GetDirection() string {
	return e.Direction
}

func (e *Entity) SetDirectionToTargetPosition(x, y int) {
	e.SetDirection(e.GetDirectionToTargetPosition(x, y))
}

func (e *Entity) GetDirectionToTargetPosition(x, y int) string {
	ax, ay := e.GetPosition()

	dx := x - ax
	dy := y - ay

	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		if dx < 0 {
			return "left"
		} else {
			return "right"
		}
	} else {
		if dy < 0 {
			return "up"
		} else {
			return "down"
		}
	}
}