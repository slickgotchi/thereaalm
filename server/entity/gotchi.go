package entity

import (
	"log"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Gotchi struct {
    Entity
	ActionSequence
	types.Inventory
	types.Stats
	GotchiId string
}

func NewGotchi(zoneId, x, y int) *Gotchi {
	// add item holder
	newItemHolder := types.NewInventory()

	// add some stats
	newStats := types.NewStats()
	newStats.AddDynamicStat("hp", 400, 400)
	newStats.AddStaticStat("attack", 5)
	newStats.AddStaticStat("harvest_duration_s", 10)
	newStats.AddStaticStat("trade_duration_s", 10)

	// make new gotchi
	return &Gotchi{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "gotchi",
			X: x,
			Y: y,
        },
        ActionSequence: ActionSequence{
			Actions: make([]types.IAction, 0),
		},
		Inventory: *newItemHolder,
		Stats: *newStats,
		GotchiId: "4285",
    }
}

func (g *Gotchi) GetSnapshotData() interface{} {
	hp, _ := g.Stats.GetStatValue("hp")
	maxHP, _ := g.Stats.GetStatMaxValue("hp")
	return struct {
		GotchiID  string `json:"gotchiId"`
		MaxHP     int    `json:"maxHp"`
		CurrentHP int    `json:"currentHp"`
	}{
		GotchiID:  g.GotchiId,
		MaxHP:     maxHP,
		CurrentHP: hp,
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
