package entity

import (
	"thereaalm/entity/entitystate"
	"thereaalm/stattypes"

	"github.com/google/uuid"
)



type Altar struct {
	Entity
	stattypes.Stats
	entitystate.State
}

func NewAltar(zoneId, x, y int) *Altar {
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
		Stats: e.StatMap,
		State: e.State,
	}
}

func (e *Altar) Update(dt_s float64) {
	pulse := e.GetStat(stattypes.Pulse)

	if pulse <= 0 {
		e.State = entitystate.Dead
	} 
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

// GetMaxPulse() already part of IMaintainable