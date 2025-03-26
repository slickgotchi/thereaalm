package action

import (
	"log"
	"thereaalm/stats"
	"thereaalm/types"
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

func (action *HarvestAction) Start() {
	action.Timer_s = action.Duration_s
}

func (action *HarvestAction) CanBeExecuted() bool {
	harvestable, _ := action.Target.(types.IHarvestable); 
	itemHolder, _ := action.Actor.(types.IInventory);
	if itemHolder == nil || harvestable == nil {
		log.Printf("Invalid actor or target in HarvestAction CanBeExecuted()")
		return false	// action is complete we have invalid actor or target
	}

	return harvestable.CanBeHarvested()
}

func (action *HarvestAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	harvestable, _ := action.Target.(types.IHarvestable); 
	itemHolder, _ := action.Actor.(types.IInventory);
	if itemHolder == nil || harvestable == nil {
		log.Printf("Invalid actor or target in HarvestAction Update()")
		return true	// action is complete we have invalid actor or target
	}

	// move to target
	tx, ty := action.Target.GetPosition()
	action.Actor.SetPosition(tx, ty +1)

	// check duration expired
	action.Timer_s -= dt_s
	if action.Timer_s <= 0 {
		typeRemoved, amountRemoved := harvestable.Harvest()
		log.Println(typeRemoved, amountRemoved)
		if typeRemoved != "" && amountRemoved > 0 {
			itemHolder.AddItem(typeRemoved, amountRemoved)
			log.Printf("%s added %d %s to inventory", action.Actor.GetType(), amountRemoved, typeRemoved)
		}

		itemHolder.DisplayInventory()

		// harvesting is complete so we return TRUE
		return true
	}

	// harvesting is not complete so we return FALSE
	return false
}