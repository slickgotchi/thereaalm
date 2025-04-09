package entity

import (
	"math/rand"
	"thereaalm/action"
	"thereaalm/entity/entitystate"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"
	"time"

	"github.com/google/uuid"
)



type LickVoid struct {
	Entity
	Stats stattypes.Stats
	entitystate.State
	SpawnInterval_s float64
	LastSpawnTime time.Duration
	MaxAliveSpawns  int
	SpawnedLicks    []interfaces.IEntity
}

func NewLickVoid(x, y int) *LickVoid {
	newStats := stattypes.NewStats()
	newStats.SetStat(stattypes.Pulse, 1000)
	newStats.SetStat(stattypes.MaxPulse, 1000)

	return &LickVoid{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "lickvoid",
			X: x,
			Y: y,
        },
		Stats: *newStats,
		State: entitystate.Active,
		LastSpawnTime: 0,
		SpawnInterval_s: 5,
		MaxAliveSpawns:  5,
		SpawnedLicks:    []interfaces.IEntity{},
    }
}

func (e *LickVoid) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Stats interface{} `json:"stats"`
		State entitystate.State `json:"state"`
	}{
		Name: "Lick Void",
		Description: "A nefarious portal lickquidators use to enter The Reaalm",
		Stats: e.Stats.StatMap,
		State: e.State,
	}
}

func (e *LickVoid) Update(dt_s float64) {
	if e.State != entitystate.Active {
		return
	}

	// Clean up removed Lickquidators (nil or no longer in a zone)
	filtered := e.SpawnedLicks[:0]
	for _, l := range e.SpawnedLicks {
		if l != nil && l.GetZone() != nil {
			filtered = append(filtered, l)
		}
	}
	e.SpawnedLicks = filtered

	// Enforce alive spawn limit
	if len(e.SpawnedLicks) >= e.MaxAliveSpawns {
		return
	}

	// Spawn if interval has passed
	if e.WorldManager.Since(e.LastSpawnTime) >= time.Duration(e.SpawnInterval_s)*time.Second {
		currX, currY := e.GetPosition()

		corners := [][2]int{
			{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
		}
		offset := corners[rand.Intn(len(corners))]

		spawnX := currX + offset[0]
		spawnY := currY + offset[1]

		lick := e.generateGenericLickquidator(spawnX, spawnY)
		e.SpawnedLicks = append(e.SpawnedLicks, lick)

		e.LastSpawnTime = e.WorldManager.Now()
	}
}


func (e *LickVoid) generateGenericLickquidator(x, y int) interfaces.IEntity {
	zone := e.GetZone()
	lickquidator := NewLickquidator(x, y)
	zone.AddEntity(lickquidator)

	lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.5,
		&types.TargetSpec{
			TargetType:      "gotchi",
			TargetCriterion: "nearest",
		}))
	lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.5,
		&types.TargetSpec{
			TargetType:      "altar",
			TargetCriterion: "nearest",
		}))

	return lickquidator
}

// custom stat modification wrappers
func (e *LickVoid) SetStat(name string, value float64) {
	e.Stats.SetStat(name, value)
}

func (e *LickVoid) GetStat(name string) float64 {
	return e.Stats.GetStat(name)
}

func (e *LickVoid) DeltaStat(name string, value float64) {
	prev := e.Stats.GetStat(name)
	e.Stats.DeltaStat(name, value)
	newVal := e.Stats.GetStat(name)

	// CUSTOM HOOK: handle Pulse stats going below 0 (death)
	if (name == stattypes.Pulse) && 
		newVal <= 0 && prev > 0 {

		// set gotchi state to dead
		e.State = entitystate.Dead
	}
}

/*
// IMaintainable functions
func (e *LickVoid) Maintain(pulseRestored int) {
	e.Stats.DeltaStat(stattypes.Pulse, pulseRestored)
	if e.Stats.GetStat(stattypes.Pulse) > e.GetStat(stattypes.MaxPulse) {
		e.Stats.SetStat(stattypes.Pulse, e.GetStat(stattypes.MaxPulse))
	}
}

func (e *LickVoid) CanBeMaintained() bool {
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
func (e *LickVoid) Rebuild(pulseRestored int) {
	e.Stats.DeltaStat(stattypes.Pulse, pulseRestored)

	// if pulse >= max pulse set our state back to active
	if e.Stats.GetStat(stattypes.Pulse) >= e.Stats.GetStat(stattypes.MaxPulse) {
		e.State = entitystate.Active
		e.Stats.SetStat(stattypes.Pulse, e.GetStat(stattypes.MaxPulse))
	}
}

func (e *LickVoid) CanBeRebuilt() bool {
	// structure must be in Dead state to be rebuilt
	if e.State != entitystate.Dead {
		return false
	}

	// ok! we can be rebuilt
	return true
}
*/