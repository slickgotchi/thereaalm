package action

import (
	"fmt"
	"log"
	"thereaalm/stats"
	"thereaalm/types"
	"time"
)

type HarvestAction struct {
	Action
	Timer_s float64
	Duration_s float64
}

func NewHarvestAction(actor, target types.IEntity, weighting float64) *HarvestAction {
	actorItemHolder, _ := actor.(types.IInventory)
	if actorItemHolder == nil {
		log.Println("failed test")
	}
	actorStats, _ := actor.(stats.IStats)
	if actorStats == nil {
		log.Println("ERROR: Harvesting actor does not have IStats, returning...")
		return nil
	}

	harvestDuration_s := actorStats.GetStat(stats.HarvestDuration_s)
	if harvestDuration_s <= 0 {
		log.Println("ERROR: Harvesting actor must have 'harvest_duration_s' stat, returning...")
		return nil
	}

	return &HarvestAction{
		Action: Action{
			Type: "harvest",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
		Timer_s: float64(harvestDuration_s),
		Duration_s: float64(harvestDuration_s),
	}
}

func (a *HarvestAction) CanBeExecuted() bool {
	harvestable, _ := a.Target.(types.IHarvestable); 
	itemHolder, _ := a.Actor.(types.IInventory);

	// actor and target of correct types?
	if itemHolder == nil || harvestable == nil {
		log.Printf("Invalid actor or target in HarvestAction CanBeExecuted()")
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(a.Target) {
		return false
	}

	// is harvestable?
	if !harvestable.CanBeHarvested() {
		return false
	}

	log.Println("Its HARVEST time!!!")
	// ok can execute
	return true
}

func (a *HarvestAction) Start() {
	a.Timer_s = a.Duration_s

	// move to target
	a.TryMoveToTargetEntity(a.Target)
}

func (a *HarvestAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	harvestable, _ := a.Target.(types.IHarvestable); 
	itemHolder, _ := a.Actor.(types.IInventory);
	if itemHolder == nil || harvestable == nil {
		log.Printf("Invalid actor or target in HarvestAction Update()")
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	a.Timer_s -= dt_s
	if a.Timer_s <= 0 {
		typeRemoved, amountRemoved := harvestable.Harvest()
		log.Println(typeRemoved, amountRemoved)
		if typeRemoved != "" && amountRemoved > 0 {
			itemHolder.AddItem(typeRemoved, amountRemoved)
			log.Printf("%s added %d %s to inventory", a.Actor.GetType(), amountRemoved, typeRemoved)

			// see if harvester has an activity log
			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				entry := types.ActivityLogEntry{
					Description: fmt.Sprintf("Harvested %d %s", amountRemoved, typeRemoved),
					LogTime: time.Now(),
				}
				activityLog.NewLogEntry(entry)
			}
		}

		itemHolder.DisplayInventory()

		// harvesting is complete so we return TRUE
		return true
	}

	// harvesting is not complete so we return FALSE
	return false
}