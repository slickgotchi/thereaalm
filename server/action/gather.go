package action

import (
	"fmt"
	"log"
	"thereaalm/jobs"
	"thereaalm/mathext"
	"thereaalm/stats"
	"thereaalm/types"
	"time"
)

type GatherAction struct {
	Action
	Timer_s float64
	Duration_s float64
}

func NewGatherAction(actor, target types.IEntity, weighting float64) *GatherAction {
	actorItemHolder, _ := actor.(types.IInventory)
	if actorItemHolder == nil {
		log.Println("failed test")
	}
	actorStats, _ := actor.(stats.IStats)
	if actorStats == nil {
		log.Println("ERROR: Harvesting actor does not have IStats, returning...")
		return nil
	}

	// ecto determines gather duration
	actorEcto := actorStats.GetStat(stats.Ecto)
	if actorEcto < 0 {
		log.Printf("Actor does not have an ecto stat. One must be assigned to use 'gather'")
		return nil
	}

	// find ecto delta from farmer peak and clamp it between 0 and 500
	deltaToPeakEcto := mathext.Abs(actorEcto - jobs.Farmer.Peak.Ecto)
	deltaToPeakEcto = mathext.Clamp(deltaToPeakEcto, 0, 500)

	// vary gather duration between 5 - 30 seconds
	alpha := float64(deltaToPeakEcto) / 500.0
	gatherDuration_s := int(5 + 25 * alpha)

	return &GatherAction{
		Action: Action{
			Type: "gather",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
		Timer_s: float64(gatherDuration_s),
		Duration_s: float64(gatherDuration_s),
	}
}

func (a *GatherAction) CanBeExecuted() bool {
	gatherable, _ := a.Target.(types.IGatherable); 
	itemHolder, _ := a.Actor.(types.IInventory);

	// actor and target of correct types?
	if itemHolder == nil || gatherable == nil {
		log.Printf("Invalid actor or target in HarvestAction CanBeExecuted()")
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(a.Target) {
		return false
	}

	// is harvestable?
	if !gatherable.CanBeGathered() {
		return false
	}

	log.Println("Its HARVEST time!!!")
	// ok can execute
	return true
}

func (a *GatherAction) Start() {
	a.Timer_s = a.Duration_s

	// move to target
	a.TryMoveToTargetEntity(a.Target)
}

func (a *GatherAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	gatherable, _ := a.Target.(types.IGatherable); 
	itemHolder, _ := a.Actor.(types.IInventory);
	if itemHolder == nil || gatherable == nil {
		log.Printf("Invalid actor or target in HarvestAction Update()")
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	a.Timer_s -= dt_s
	if a.Timer_s <= 0 {
		typeRemoved, amountRemoved := gatherable.Gather()
		// log.Println(typeRemoved, amountRemoved)
		if typeRemoved != "" && amountRemoved > 0 {
			// add item to item holder
			itemHolder.AddItem(typeRemoved, amountRemoved)

			// remove spark equivalent to the gather duration from the gatherer
			if gathererStats, ok := a.Actor.(stats.IStats); ok {
				gathererStats.DeltaStat(stats.Spark, -int(a.Duration_s))
			}

			// see if gatherer has an activity log
			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				entry := types.ActivityLogEntry{
					Description: fmt.Sprintf("Gathered %d %s", amountRemoved, typeRemoved),
					LogTime: time.Now(),
				}
				activityLog.NewLogEntry(entry)
			}
		}

		// harvesting is complete so we return TRUE
		return true
	}

	// harvesting is not complete so we return FALSE
	return false
}