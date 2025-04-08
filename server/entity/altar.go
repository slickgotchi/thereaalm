package entity

import (
	"thereaalm/entity/entitystate"
	"thereaalm/stattypes"

	"github.com/google/uuid"
)



type Altar struct {
	Entity
	Stats stattypes.Stats
	entitystate.State
	BuffRange int // Add range field
}

func NewAltar(x, y int) *Altar {
	newStats := stattypes.NewStats()
	newStats.SetStat(stattypes.Pulse, 1000)
	newStats.SetStat(stattypes.MaxPulse, 1000)

	return &Altar{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "altar",
			X: x,
			Y: y,
        },
		Stats: *newStats,
		State: entitystate.Active,
		BuffRange: 10,
    }
}

func (e *Altar) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Stats interface{} `json:"stats"`
		State entitystate.State `json:"state"`
	}{
		Name: "Gotchi Altar",
		Description: "While active, imbues nearby gotchis with action duration bonuses",
		Stats: e.Stats.StatMap,
		State: e.State,
	}
}

func (e *Altar) Update(dt_s float64) {

}

// IMaintainable functions
func (e *Altar) Maintain(pulseRestored int) {
	e.Stats.DeltaStat(stattypes.Pulse, pulseRestored)
	if e.Stats.GetStat(stattypes.Pulse) > e.GetStat(stattypes.MaxPulse) {
		e.Stats.SetStat(stattypes.Pulse, e.GetStat(stattypes.MaxPulse))
	}
}

func (e *Altar) CanBeMaintained() bool {
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
func (e *Altar) Rebuild(pulseRestored int) {
	e.Stats.DeltaStat(stattypes.Pulse, pulseRestored)

	// if pulse >= max pulse set our state back to active
	if e.Stats.GetStat(stattypes.Pulse) >= e.Stats.GetStat(stattypes.MaxPulse) {
		e.State = entitystate.Active
		e.Stats.SetStat(stattypes.Pulse, e.GetStat(stattypes.MaxPulse))
	}
}

func (e *Altar) CanBeRebuilt() bool {
	// structure must be in Dead state to be rebuilt
	if e.State != entitystate.Dead {
		return false
	}

	// ok! we can be rebuilt
	return true
}

func (e *Altar) GetBuffRange() int {
    return e.BuffRange
}

func (e *Altar) GetSpeedBuffMultiplier() float64 {
    return 1.2 // 20% speed increase when in range
}

func (e *Altar) IsBuffActive() bool {
    return e.State == entitystate.Active && e.GetStat(stattypes.Pulse) > 0
}

// custom stat modification wrappers
func (e *Altar) SetStat(name string, value int) {
	e.Stats.SetStat(name, value)
}

func (e *Altar) GetStat(name string) int {
	return e.Stats.GetStat(name)
}

func (e *Altar) DeltaStat(name string, value int) {
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