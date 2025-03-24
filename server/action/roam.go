package action

import (
	"log"
	"thereaalm/types"
	"time"
)

type RoamAction struct {
	Type      string
	Weighting float64
	Actor     types.IEntity
	Target    types.IEntity

	// Additional fields
	NewX, NewY  int
	Completed   bool
	StartTime   time.Time
	Duration    time.Duration
}

func NewRoamAction(actor types.IEntity, weighting float64) *RoamAction {
	return &RoamAction{
		Type:      "Roam",
		Weighting: weighting,
		Actor:     actor,
		Completed: false,
		Duration:  5 * time.Second, // Set the duration to 5 seconds
		StartTime: time.Now(),
	}
}

// CanBeExecuted always returns true for RoamAction
func (r *RoamAction) CanBeExecuted() bool {
	return true
}

// Update moves the actor to the new location and completes the action after 5 seconds
func (r *RoamAction) Update(dt_s float64) bool {
	// If the action has already been completed, return true
	if r.Completed {
		return true
	}

	// Check if 5 seconds have passed
	if time.Since(r.StartTime) >= r.Duration {
		r.Completed = true
		log.Printf("Actor %s completed roaming action.\n", r.Actor.GetUUID())
		return true
	}

	// cast actor to a IZoned type
	zonedActor, _ := r.Actor.(types.IZoned)
	if zonedActor == nil {
		return true
	}

	// If not completed, attempt to find a new empty cell using the zone's FindNearbyEmptyCell method
	zone := zonedActor.GetZone() // Get the actor's zone
	actorX, actorY := r.Actor.GetPosition()

	// Use the zone's FindNearbyEmptyCell with radius 3
	newX, newY, found := zone.FindNearbyEmptyCell(actorX, actorY, 3)
	if found {
		r.NewX, r.NewY = newX, newY
		// Move the actor to the new position
		r.Actor.SetPosition(newX, newY)
		log.Printf("Actor %s moved to (%d, %d)\n", r.Actor.GetUUID(), newX, newY)
	}

	return false
}
