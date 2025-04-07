package action

import (
	"log"
	"math/rand"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"
	"time"
)

type RoamAction struct {
	Action
	StartTime   time.Time
	Duration    time.Duration
}

func NewRoamAction(actor interfaces.IEntity, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *RoamAction {

	a := &RoamAction{
		Action: Action{
			Type:      "roam",
			Weighting: weighting,
			Actor:     actor,
			Target: nil,
		},
		Duration:  5 * time.Second, // Set the duration to 5 seconds
		StartTime: time.Now(),
	}

	a.SetFallbackTargetSpec(fallbackTargetSpec)

	return a
}

func (r *RoamAction) IsValidTarget(potentialTarget interfaces.IEntity) bool {
	return true
}

func (r *RoamAction) IsValidActor(potentialActor interfaces.IEntity) bool {
	return true
}

func (r *RoamAction) Start() {
    r.StartTime = time.Now()
    r.Duration = time.Duration(5+rand.Float64()*(15-5)) * time.Second

	// attempt to find a new empty cell using the zone's FindNearbyEmptyCell method
	zone := r.Actor.GetZone() // Get the actor's zone
	actorX, actorY := r.Actor.GetPosition()

	// use ecto to govern roam radius (between 2 - 10)
	actorStats, _ := r.Actor.(interfaces.IStats)
	if actorStats == nil {
		log.Printf("Roam assigned to actor with no stats!")
		return
	}

	// find delta from peak explorer ecto
	actorEcto := actorStats.GetStat(stattypes.Ecto)

	alpha := 1.0 - actorEcto / 1000
	explorationRadius := 2 + int(alpha * 8.0)

	// Use the zone's FindNearbyEmptyCell with radius 3
	newX, newY, found := zone.FindNearbyEmptyTile(actorX, actorY, explorationRadius, 1)
	if found {
		// set direction to new position
		r.Actor.SetDirectionToTargetPosition(newX, newY)

		// Move the actor to the new position
		r.Actor.SetPosition(newX, newY)

		// reduce pulse (our "stability")
		actorStats.DeltaStat(stattypes.Pulse, -5)
	}
}



// Update moves the actor to the new location and completes the action after 5 seconds
func (r *RoamAction) Update(dt_s float64) bool {
	// Check if 5 seconds have passed
	if time.Since(r.StartTime) >= r.Duration {
		// log.Printf("Actor %s completed roaming action.\n", r.Actor.GetUUID())
		return true
	}

	return false
}
