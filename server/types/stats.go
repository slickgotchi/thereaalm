package types

type StatType int

const (
	StatTypeStatic = iota
	StatTypeDynamic
)

type Stat struct {
	Type StatType
	Value int
	MaxValue int
}

type IStats interface {
    AddStaticStat(name string, value int)
	AddDynamicStat(name string, value, maxValue int)
	GetValue(name string) (int, bool)
	GetMaxValue(name string) (int, bool)
	SetValue(name string, value int) bool
	SetMaxValue(name string, max int) bool
	DeltaValue(name string, delta int)
}

type Stats struct {
	Entries map[string]Stat
}

func NewStats() *Stats {
	return &Stats{
		Entries: make(map[string]Stat),
	}
}

// add a new static stat
func (s * Stats) AddStaticStat(name string, value int) {
	s.Entries[name] = Stat{
		Type: StatTypeStatic,
		Value: value,
	}
}

// add a new dynamic stat
func (s * Stats) AddDynamicStat(name string, value, maxValue int) {
	s.Entries[name] = Stat{
		Type: StatTypeDynamic,
		Value: value,
		MaxValue: maxValue,
	}
}

// GetValue returns the value of a stat (for static stats) or current value (for dynamic stats)
func (s *Stats) GetValue(name string) (int, bool) {
    stat, exists := s.Entries[name]
    if !exists {
        return 0, false
    }
    return stat.Value, true
}

// GetMaxValue returns the max value of a dynamic stat (returns 0 for static stats)
func (s *Stats) GetMaxValue(name string) (int, bool) {
    stat, exists := s.Entries[name]
    if !exists || stat.Type != StatTypeDynamic {
        return 0, false
    }
    return stat.MaxValue, true
}

// SetValue sets the value of a stat (for static stats) or current value (for dynamic stats)
func (s *Stats) SetValue(name string, value int) bool {
    stat, exists := s.Entries[name]
    if !exists {
        return false
    }
    stat.Value = value
    // For dynamic stats, ensure current doesn't exceed max
    if stat.Type == StatTypeDynamic && stat.Value > stat.MaxValue {
        stat.Value = stat.MaxValue
    }
    s.Entries[name] = stat
    return true
}

// SetMaxValue sets the max value of a dynamic stat
func (s *Stats) SetMaxValue(name string, max int) bool {
    stat, exists := s.Entries[name]
    if !exists || stat.Type != StatTypeDynamic {
        return false
    }
    stat.MaxValue = max
    // Adjust current value if it exceeds the new max
    if stat.Value > stat.MaxValue {
        stat.Value = stat.MaxValue
    }
    s.Entries[name] = stat
    return true
}

// DeltaValue modifies the value of a stat by a delta (for static or dynamic stats)
func (s *Stats) DeltaValue(name string, delta int) bool {
    stat, exists := s.Entries[name]
    if !exists {
        return false
    }
    stat.Value += delta
    // For dynamic stats, clamp between 0 and max
    if stat.Type == StatTypeDynamic {
        if stat.Value > stat.MaxValue {
            stat.Value = stat.MaxValue
        }
        if stat.Value < 0 {
            stat.Value = 0
        }
    }
    s.Entries[name] = stat
    return true
}