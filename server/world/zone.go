// thereaalm/world/zone.go
package world

import (
	"log"
	"thereaalm/entity"
)

type Zone struct {
    ID       int
    Entities []entity.Entity
    Width    int
    Height   int
    X        int
    Y        int
}

func NewZone(id, width, height, x, y int) *Zone {
    return &Zone{
        ID:       id,
        Entities: []entity.Entity{},
        Width:    width,
        Height:   height,
        X:        x,
        Y:        y,
    }
}

func (z *Zone) AddEntity(e entity.Entity) {
    z.Entities = append(z.Entities, e)
}

func (z *Zone) Update() {
    log.Printf("Updating Zone %d with %d entities", z.ID, len(z.Entities))
    for _, e := range z.Entities {
        e.Update()
        // pos := e.GetPosition()
        // log.Printf("Entity at (%d, %d) in Zone %d", pos.X, pos.Y, pos.ZoneID)
    }
}