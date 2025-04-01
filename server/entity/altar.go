package entity

import (
	"thereaalm/entity/entitystate"
	"thereaalm/stats"

	"github.com/google/uuid"
)



type Altar struct {
	Entity
	stats.Stats
	State entitystate.EntityState
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
		State: entitystate.Active,
		MaxPulse: 1000,
    }
}

func (e *Altar) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Stats interface{} `json:"stats"`
	}{
		Name: "Gotchi Altar",
		Description: "While active, imbues nearby gotchis with action duration bonuses",
		Stats: e.StatMap,
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
