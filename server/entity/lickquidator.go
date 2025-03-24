package entity

import (
	// "log"
	"thereaalm/action"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Lickquidator struct {
    Entity
	action.ActionPlan
	types.Inventory
	types.Stats
}

func NewLickquidator(zoneId, x, y int) *Lickquidator {
	// make them hold the "Tongue" item
	newInventory := types.NewInventory()
	newInventory.Items["tongue"] = 1

	// give a base hp stat
	newStats := types.NewStats()
	newStats.AddDynamicStat("hp", 50, 50)
	newStats.AddStaticStat("attack", 3)

    return &Lickquidator{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "lickquidator",
			X: x,
			Y: y,
        },
        ActionPlan: action.ActionPlan{
			Actions: make([]types.IAction, 0),
		},
		Inventory: *newInventory,
		Stats: *newStats,
    }
}

func (l *Lickquidator) GetSnapshotData() interface{} {
	return struct {
		Stats interface{} `json:"stats"`
	}{
		Stats: l.Stats,
	}
}

func (l *Lickquidator) Update(dt_s float64) {
	// process actions
	// l.ProcessActions(dt_s)
}
