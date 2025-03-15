package ai

// short/long term memory handling

import "time"

// GotchiMind represents the brain of a Gotchi.
type GotchiMind struct {
    Personality     map[string]int // Fixed traits
    ShortTermMemory []string       // Recent experiences
    LongTermMemory  []string       // Key life memories
    LastInteraction time.Time      // Last AI update
}

// NewGotchiMind initializes a new GotchiMind.
func NewGotchiMind() GotchiMind {
    return GotchiMind{
        Personality: map[string]int{"curiosity": 50, "patience": 50},
    }
}

// AddMemory records a new memory in short-term storage.
func (m *GotchiMind) AddMemory(memory string) {
    m.ShortTermMemory = append(m.ShortTermMemory, memory)
}

func (m *GotchiMind) Update() {
    
}