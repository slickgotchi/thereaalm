package stattypes

import (
	"log"
)

// Stat constants to avoid mistyped stat names.
const (
    Ecto = "ecto"
    Spark = "spark"
    Pulse = "pulse"
    MaxPulse = "maxpulse"
    StakedGHST = "stakedGhst"
    TreatTotal = "treatTotal"
	// HpCurrent = "hp_current"
	// HpMax     = "hp_max"
	// Attack    = "attack"
	// HarvestDuration_s   = "harvest_duration_s"
    // TradeDuration_s = "trade_duration_s"
)



// STATS
type Stats struct {
	StatMap map[string]float64
}

func NewStats() *Stats {
	return &Stats{
		StatMap: make(map[string]float64),
	}
}


// SetStat sets the value of a given stat.
func (s *Stats) SetStat(name string, value float64) {
    s.StatMap[name] = value
}

// GetStat retrieves the value of a given stat.
// If the stat does not exist, logs a warning and returns 0.
func (s *Stats) GetStat(name string) float64 {
    statValue, ok := s.StatMap[name]
    if !ok {
        log.Printf("WARNING: Stat '%s' not found.", name)
        return -1
    }
    return statValue
}

// DeltaStat modifies a stat by a given delta value.
// If the stat does not exist, logs a warning.
func (s *Stats) DeltaStat(name string, value float64) {
    _, ok := s.StatMap[name]
    if !ok {
        log.Printf("WARNING: Cannot modify unknown stat '%s'.", name)
        return
    }
    s.StatMap[name] += value
}
