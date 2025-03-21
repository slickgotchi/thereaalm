package entity

import (
	"log"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Gotchi struct {
    Entity
    Movable
	ActionSequence
	ItemHolder
}

func NewGotchi(zoneId, x, y int) *Gotchi {
    return &Gotchi{
        Entity: Entity{
            ID:   types.EntityUUID(uuid.New()),
            Type: "Gotchi",
        },
        Movable: Movable{
			ZoneID: zoneId,
            X: x,
            Y: y,
        },
        ActionSequence: ActionSequence{
			Actions: make([]types.IAction, 0),
		},
		ItemHolder: *NewItemHolder(),
    }
}

func (g *Gotchi) Update(dt_s float64) {
	log.Printf("Gotchi at (%d, %d)", g.X, g.Y)
	g.DisplayInventory()

	// process the oldest action until its done
	if len(g.Actions) > 0 {
		isComplete := g.Actions[0].Update(dt_s)

		if isComplete {
			g.Actions = g.Actions[1:]
		}
	}
}
