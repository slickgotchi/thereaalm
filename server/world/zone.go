// thereaalm/world/zone.go
package world

import (
	"thereaalm/types"
)

type Zone struct {
    ID       int
    Entities []types.IEntity
    Width    int
    Height   int
    X        int
    Y        int
}

func NewZone(id, width, height, x, y int) *Zone {
    return &Zone{
        ID:       id,
        Entities: []types.IEntity{},
        Width:    width,
        Height:   height,
        X:        x,
        Y:        y,
    }
}

func (z *Zone) AddEntity(e types.IEntity) {
    z.Entities = append(z.Entities, e)
}

func (z *Zone) Update(dt_s float64) {
    // log.Printf("Updating Zone %d with %d entities", z.ID, len(z.Entities))
    for _, e := range z.Entities {
        e.Update(dt_s)
        // pos := e.GetPosition()
        // log.Printf("Entity at (%d, %d) in Zone %d", pos.X, pos.Y, pos.ZoneID)
    }
}