package entity

import (
	// "log"
	"thereaalm/entity/entitystate"
	"thereaalm/stattypes"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Shop struct {
    Entity
	Stats stattypes.Stats
	types.Inventory
	entitystate.State
}

func NewShop(x, y int) *Shop {
	// start show with gold
	itemHolder := types.NewInventory()
	itemHolder.Items["gold"] = 10000

	newStats := stattypes.NewStats()
	newStats.SetStat(stattypes.Pulse, 1000)
	newStats.SetStat(stattypes.MaxPulse, 1000)

    return &Shop{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "shop",
            X: x,
            Y: y,
        },
        Inventory: *itemHolder,
		Stats: *newStats,
		State: entitystate.Active,
    }
}

func (s *Shop) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Inventory map[string]int `json:"inventory"`
		Stats interface{} `json:"stats"`
		State entitystate.State `json:"state"`
	}{
		Name: s.Type,
		Description: "Buy and sell items from one convenient location",
		Inventory: s.Items,
		Stats: s.Stats.StatMap,
		State: s.State,
	}
}

func (s *Shop) Update(dt_s float64) {

}

// IMaintainable functions
func (e *Shop) Maintain(pulseRestored float64) {
	e.Stats.DeltaStat(stattypes.Pulse, pulseRestored)
	if e.Stats.GetStat(stattypes.Pulse) > e.Stats.GetStat(stattypes.MaxPulse) {
		e.Stats.SetStat(stattypes.Pulse, e.Stats.GetStat(stattypes.MaxPulse))
	}
}

func (e *Shop) CanBeMaintained() bool {
	// structure must be Active state to be maintained
	if e.State != entitystate.Active {
		return false
	}

	// don't allow maintenance on structures above 80% pulse
	currPulse := e.Stats.GetStat(stattypes.Pulse)
	maxPulse := e.GetStat(stattypes.MaxPulse)
	if float64(currPulse) >= float64(maxPulse)*0.8 {
		return false
	}

	// ok! we can be maintained
	return true
}

// IRebuildable functions
func (e *Shop) Rebuild(pulseRestored float64) {
	e.Stats.DeltaStat(stattypes.Pulse, pulseRestored)

	// if pulse >= max pulse set our state back to active
	if e.Stats.GetStat(stattypes.Pulse) >= e.Stats.GetStat(stattypes.MaxPulse) {
		e.State = entitystate.Active
		e.Stats.SetStat(stattypes.Pulse, e.Stats.GetStat(stattypes.MaxPulse))
	}
}

func (e *Shop) CanBeRebuilt() bool {
	// structure must be in Dead state to be rebuilt
	if e.State != entitystate.Dead {
		return false
	}

	// ok! we can be rebuilt
	return true
}

// custom stat modification wrappers
func (e *Shop) SetStat(name string, value float64) {
	e.Stats.SetStat(name, value)
}

func (e *Shop) GetStat(name string) float64 {
	return e.Stats.GetStat(name)
}

func (e *Shop) DeltaStat(name string, value float64) {
	prev := e.Stats.GetStat(name)
	e.Stats.DeltaStat(name, value)
	newVal := e.Stats.GetStat(name)

	// CUSTOM HOOK: handle ESP stats going below 0 (death)
	if (name == stattypes.Pulse) && 
		newVal <= 0 && prev > 0 {

		// set gotchi state to dead
		e.State = entitystate.Dead
	}
}