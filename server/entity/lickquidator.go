package entity

import (
	// "log"

	"thereaalm/action"
	"thereaalm/entity/entitystate"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"
	"thereaalm/utils"

	"github.com/google/uuid"
)

type Lickquidator struct {
    Entity
	action.ActionPlan
	types.Inventory
	stattypes.Stats
	entitystate.State
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
		State: entitystate.Active,
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
	// ensure spark and ecto stats stay constant
	l.SetStat(stattypes.Ecto, 50)
	l.SetStat(stattypes.Spark, 50)

	// process actions
	l.ProcessActions(dt_s)

	// Check if the Lickquidator's Pulse is zero or less
    if l.Stats.GetStat(stattypes.Pulse) <= 0 {
        // log.Printf("Lickquidator %s has died (Pulse <= 0), removing from zone", l.GetUUID().String())
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

// custom stat modification wrappers
func (e *Lickquidator) SetStat(name string, value float64) {
	e.Stats.SetStat(name, value)
}

func (e *Lickquidator) GetStat(name string) float64 {
	return e.Stats.GetStat(name)
}

func (e *Lickquidator) DeltaStat(name string, value float64) {
	prev := e.Stats.GetStat(name)
	e.Stats.DeltaStat(name, value)
	newVal := e.Stats.GetStat(name)

	// ensure ecto, spark and pulse stay within range
	e.Stats.SetStat(stattypes.Ecto, utils.Clamp(e.Stats.GetStat(stattypes.Ecto), 0, 1000))
	e.Stats.SetStat(stattypes.Spark, utils.Clamp(e.Stats.GetStat(stattypes.Spark), 0, 1000))
	e.Stats.SetStat(stattypes.Pulse, utils.Clamp(e.Stats.GetStat(stattypes.Pulse), 0, 1000))

	// CUSTOM HOOK: handle pulse goes to zero
	if name == stattypes.Pulse && 
		newVal <= 0 && prev > 0 {

		// set gotchi state to dead
		e.State = entitystate.Dead

		// log.Printf("Lickquidator %s has died (Pulse <= 0), removing from zone", e.GetUUID().String())
        // Remove the Lickquidator from its zone
        zone := e.GetZone()
        if zone != nil {
            zone.RemoveEntity(e)
        }
        // Clean up the ActionPlan to prevent memory leaks
        e.ActionPlan.Actions = nil
        e.ActionPlan.CurrentAction = nil
	}
}
