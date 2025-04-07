package entity

import (
	"thereaalm/entity/entitystate"
	"thereaalm/stats"

	"github.com/google/uuid"
)



type Altar struct {
	Entity
	stats.Stats
	entitystate.State
	MaxPulse int
}

func NewAltar(zoneId, x, y int) *Altar {
	newStats := stats.NewStats()
	newStats.SetStat(stats.Pulse, 1000)
	newStats.SetStat(stats.MaxPulse, 1000)

	return &Altar{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "altar",
			X: x,
			Y: y,
        },
		Stats: *newStats,
		State: entitystate.Active,
		MaxPulse: 1000,
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
	pulse := e.GetStat(stats.Pulse)

	if pulse <= 0 {
		e.State = entitystate.Dead
	} 
	
	if e.State == entitystate.Dead {
		// check if it can be made active again
		if pulse >= e.MaxPulse {
			e.State = entitystate.Active
		}
	} else if e.State == entitystate.Active {
		// do active altar stuff
	} 
}

// IMaintainable functions
func (e *Altar) Maintain(pulseRestored int) {
	e.Stats.DeltaStat(stats.Pulse, pulseRestored)
	if e.Stats.GetStat(stats.Pulse) > e.MaxPulse {
		e.Stats.SetStat(stats.Pulse, e.MaxPulse)
	}
}

func (e *Altar) CanBeMaintained() bool {
	// structure must be Active state to be maintained
	if e.State != entitystate.Active {
		return false
	}

	// don't allow maintenance on structures above 80% pulse
	currPulse := e.Stats.GetStat(stats.Pulse)
	maxPulse := e.MaxPulse
	if float64(currPulse) >= float64(maxPulse)*0.8 {
		return false
	}

	// ok! we can be maintained
	return true
}

func (e *Altar) GetMaxPulse() int {
	return e.MaxPulse
}

// IRebuildable functions
func (e *Altar) Rebuild(pulseRestored int) {
	e.Stats.DeltaStat(stats.Pulse, pulseRestored)
	if e.Stats.GetStat(stats.Pulse) > e.MaxPulse {
		e.Stats.SetStat(stats.Pulse, e.MaxPulse)
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