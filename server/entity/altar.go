package entity

import (
	"thereaalm/entity/state"
	"thereaalm/stats"

	"github.com/google/uuid"
)



type Altar struct {
	Entity
	stats.Stats
	State state.EntityState
	MaxPulse int
}

func NewAltar(zoneId, x, y int) *Altar {
	newStats := stats.NewStats()
	newStats.SetStat(stats.Pulse, 1000)

	return &Altar{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "altar",
			X: x,
			Y: y,
        },
		Stats: *newStats,
		State: state.Active,
		MaxPulse: 1000,
    }
}

func (e *Altar) GetSnapshotData() interface{} {
	return struct {
		Stats interface{} `json:"stats"`
	}{
		Stats: e.StatMap,
	}
}

func (e *Altar) Update(dt_s float64) {
	pulse := e.GetStat(stats.Pulse)

	if pulse <= 0 {
		e.State = state.Dead
	} 
	
	if e.State == state.Dead {
		// check if it can be made active again
		if pulse >= e.MaxPulse {
			e.State = state.Active
		}
	} else if e.State == state.Active {
		// do active altar stuff
	} 
}
