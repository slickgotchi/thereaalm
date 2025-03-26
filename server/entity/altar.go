package entity

import (
	"thereaalm/stats"

	"github.com/google/uuid"
)
type Altar struct {
	Entity
	stats.Stats
}

func NewAltar(zoneId, x, y int) *Altar {
	newStats := stats.NewStats()
	newStats.SetStat(stats.HpCurrent, 1000)
	newStats.SetStat(stats.HpMax, 1000)

	return &Altar{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "altar",
			X: x,
			Y: y,
        },
		Stats: *newStats,
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

}
