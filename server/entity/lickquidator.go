package entity

import (
	// "log"
	"log"
	"thereaalm/action"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Lickquidator struct {
    Entity
	action.ActionPlan
	types.Inventory
	stattypes.Stats
}

func NewLickquidator(x, y int) *Lickquidator {
	// make them hold the "Tongue" item
	newInventory := types.NewInventory()
	newInventory.Items["tongue"] = 1

	// give a base hp stat
	newStats := stattypes.NewStats()
	newStats.SetStat(stattypes.Ecto, 50)
	newStats.SetStat(stattypes.Spark, 50)
	newStats.SetStat(stattypes.Pulse, 50)
	newStats.SetStat(stattypes.MaxPulse, 50)

    return &Lickquidator{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "lickquidator",
			X: x,
			Y: y,
        },
        ActionPlan: action.ActionPlan{
			Actions: make([]interfaces.IAction, 0),
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
		Direction string `json:"direction"`
	}{
		Name: l.Type,
		Description: "The arch enemies of the Gotchi-kin, born from the souls of liquidated traders",
		Stats: l.Stats.StatMap,
		Direction: l.Direction,
	}
}

func (l *Lickquidator) Update(dt_s float64) {
	// process actions
	l.ProcessActions(dt_s)

	// Check if the Lickquidator's Pulse is zero or less
    if l.Stats.GetStat(stattypes.Pulse) <= 0 {
        log.Printf("Lickquidator %s has died (Pulse <= 0), removing from zone", l.GetUUID().String())
        // Remove the Lickquidator from its zone
        zone := l.GetZone()
        if zone != nil {
            zone.RemoveEntity(l)
        }
        // Clean up the ActionPlan to prevent memory leaks
        l.ActionPlan.Actions = nil
        l.ActionPlan.CurrentAction = nil
        return
    }
}
