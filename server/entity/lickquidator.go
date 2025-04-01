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
	newStats.SetStat(stats.Ecto, 50)
	newStats.SetStat(stats.Spark, 50)
	newStats.SetStat(stats.Pulse, 50)

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
		Name string `json:"name"`
		Description string `json:"description"`
		Stats interface{} `json:"stats"`
	}{
		Name: l.Type,
		Description: "The arch enemies of the Gotchi-kin, born from the souls of liquidated traders",
		Stats: l.Stats.StatMap,
	}
}

func (l *Lickquidator) Update(dt_s float64) {
	// process actions
	l.ProcessActions(dt_s)
}
