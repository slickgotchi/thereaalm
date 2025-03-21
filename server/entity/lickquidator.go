package entity

import (
	"log"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Lickquidator struct {
    Entity
    Movable
	ActionSequence
	ItemHolder
}

func NewLickquidator(zoneId, x, y int) *Lickquidator {
	// make them hold the "Tongue" item
	newInventory := NewItemHolder()
	newInventory.Items["tongue"] = 1

    return &Lickquidator{
        Entity: Entity{
            ID:   types.EntityUUID(uuid.New()),
            Type: "lickquidator",
        },
        Movable: Movable{
			ZoneID: zoneId,
            X: x,
            Y: y,
        },
        ActionSequence: ActionSequence{
			Actions: make([]types.IAction, 0),
		},
		ItemHolder: *newInventory,
    }
}

func (g *Lickquidator) Update(dt_s float64) {
	log.Printf("Lickquidator at (%d, %d)", g.X, g.Y)
	g.DisplayInventory()

	// process the oldest action until its done
	if len(g.Actions) > 0 {
		isComplete := g.Actions[0].Update(dt_s)

		if isComplete {
			g.Actions = g.Actions[1:]
		}
	}
}
