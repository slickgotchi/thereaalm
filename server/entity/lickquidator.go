package entity

import (
	// "log"
	"thereaalm/action"
	"thereaalm/stats"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Lickquidator struct {
    Entity
	action.ActionPlan
	types.Inventory
	stats.Stats
}

func NewLickquidator(zoneId, x, y int) *Lickquidator {
	// make them hold the "Tongue" item
	newInventory := types.NewInventory()
	newInventory.Items["tongue"] = 1

	// give a base hp stat
	newStats := stats.NewStats()
	newStats.SetStat(stats.HpCurrent, 50)
	newStats.SetStat(stats.HpMax, 50)
	newStats.SetStat(stats.Attack, 3)

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
		Stats: l.Stats.StatMap,
	}
}

func (l *Lickquidator) Update(dt_s float64) {
	// process actions
	l.ProcessActions(dt_s)
}
