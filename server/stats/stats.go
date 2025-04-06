package stats

import (
	"log"
)

// Stat constants to avoid mistyped stat names.
const (
    Ecto = "ecto"
    Spark = "spark"
    Pulse = "pulse"
    MaxPulse = "maxpulse"
	// HpCurrent = "hp_current"
	// HpMax     = "hp_max"
	// Attack    = "attack"
	// HarvestDuration_s   = "harvest_duration_s"
    // TradeDuration_s = "trade_duration_s"
)

type IStats interface {
    SetStat(name string, value int)
    GetStat(name string) int
    DeltaStat(name string, value int)
}

// STATS
type Stats struct {
	StatMap map[string]int
}

func NewStats() *Stats {
	return &Stats{
		StatMap: make(map[string]int),
	}
}


// SetStat sets the value of a given stat.
func (s *Stats) SetStat(name string, value int) {
    s.StatMap[name] = value
}

// GetStat retrieves the value of a given stat.
// If the stat does not exist, logs a warning and returns 0.
func (s *Stats) GetStat(name string) int {
    statValue, ok := s.StatMap[name]
    if !ok {
        log.Printf("WARNING: Stat '%s' not found.", name)
        return -1
    }
    return statValue
}

// DeltaStat modifies a stat by a given delta value.
// If the stat does not exist, logs a warning.
func (s *Stats) DeltaStat(name string, value int) {
    _, ok := s.StatMap[name]
    if !ok {
        log.Printf("WARNING: Cannot modify unknown stat '%s'.", name)
        return
    }
    s.StatMap[name] += value
}
