package interfaces

import "time"

type IBuffConsumer interface {
    IEntity                    // Embed IEntity for position and zone access
    GetEffectiveSpeedMultiplier() float64 // Combined speed multiplier (base + buffs)
    SetBuffMultiplier(multiplier float64) // Set the buff-specific multiplier
    GetLastBuffCheck() time.Duration      // When buffs were last checked
    SetLastBuffCheck(t time.Duration)     // Update last check time
}

type IBuffProvider interface {
    GetBuffRange() int           // Range in tiles
    GetSpeedBuffMultiplier() float64 // Speed buff (e.g., 1.2)
    IsBuffActive() bool          // Whether buff is active
}