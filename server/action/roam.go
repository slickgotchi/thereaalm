package action

import (
	"log"
	"thereaalm/types"
	"time"
)

type RoamAction struct {
	Action

	// Additional fields
	// NewX, NewY  int
	StartTime   time.Time
	Duration    time.Duration
}

func NewRoamAction(actor types.IEntity, weighting float64) *RoamAction {
	return &RoamAction{
		Action: Action{
			Type:      "Roam",
			Weighting: weighting,
			Actor:     actor,
			Target: nil,
		},
		Duration:  5 * time.Second, // Set the duration to 5 seconds
		StartTime: time.Now(),
	}
}

func (r *RoamAction) Start() {
	r.StartTime = time.Now()
}

// CanBeExecuted always returns true for RoamAction
func (r *RoamAction) CanBeExecuted() bool {
	return true
}

// Update moves the actor to the new location and completes the action after 5 seconds
func (r *RoamAction) Update(dt_s float64) bool {
	// Check if 5 seconds have passed
	if time.Since(r.StartTime) >= r.Duration {
		log.Printf("Actor %s completed roaming action.\n", r.Actor.GetUUID())
		return true
	}

	// If not completed, attempt to find a new empty cell using the zone's FindNearbyEmptyCell method
	zone := r.Actor.GetZone() // Get the actor's zone
	actorX, actorY := r.Actor.GetPosition()

	// Use the zone's FindNearbyEmptyCell with radius 3
	newX, newY, found := zone.FindNearbyEmptyTile(actorX, actorY, 3)
	if found {
		// set direction to new position
		r.Actor.SetDirectionToTargetPosition(newX, newY)

		// Move the actor to the new position
		r.Actor.SetPosition(newX, newY)
		log.Printf("Actor %s moved to (%d, %d)\n", r.Actor.GetUUID(), newX, newY)
	}

	return false
}
