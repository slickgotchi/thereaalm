package explorationactions

import (
	"log"
	"thereaalm/action"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"
	"thereaalm/utils"
)

type RoamAction struct {
	action.Action
	
	Duration_s    float64
	Timer_s float64
}

func NewRoamAction(actor interfaces.IEntity, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *RoamAction {

	wm := actor.GetZone().GetWorldManager()

	// check for gotchi job multiplier
	jobMultiplier, err := utils.GetJobActionMultiplier(actor, "roam")
	if err != nil {
		log.Printf("ERROR [%s]: Invalid actor or action name, returning...", utils.GetFuncName())
		jobMultiplier = 1
	}

	// roam duration between 100 - 300 seconds
	roamDuration_s := 50 + 250 / jobMultiplier

	a := &RoamAction{
		Action: action.Action{
			Type:      "roam",
			Weighting: weighting,
			Actor:     actor,
			Target: nil,
			WorldManager: wm,
		},
		Duration_s:  roamDuration_s, // Set the duration to 5 seconds
		Timer_s: roamDuration_s,
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
	// attempt to find a new empty cell using the zone's FindNearbyEmptyCell method
	// zone := r.Actor.GetZone() // Get the actor's zone
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

	if r.WorldManager == nil {
		log.Println("Error: We do not have a valid WorldManager")
		return
	}

	// Use the zone's FindNearbyEmptyCell with radius 3
	newX, newY, found := r.WorldManager.FindNearbyAvailablePosition(actorX, actorY, explorationRadius, 1)
	if found {
		// set direction to new position
		r.Actor.SetDirectionToTargetPosition(newX, newY)

		// Move the actor to the new position
		r.Actor.SetPosition(newX, newY)

		// reduce pulse (our "stability")
		actorStats.DeltaStat(stattypes.Pulse, -0.5)
	}
}



// Update moves the actor to the new location and completes the action after 5 seconds
func (r *RoamAction) Update(dt_s float64) bool {
	r.Timer_s -= dt_s
	if r.Timer_s <= 0 {
		// log.Printf("Actor %s completed roaming action.\n", r.Actor.GetUUID())
		return true
	}

	return false
}
