package action

import (
	"log"
	"thereaalm/types"
)

type HarvestAction struct {
	Action
	Timer_s float64
	Duration_s float64
}

func NewHarvestAction(actor, target types.IEntity, duration_s float64) *HarvestAction {
	return &HarvestAction{
		Action: Action{
			Type: "Harvest",
			IsActive: false,
			Actor: actor,
			Target: target,
		},
		Timer_s: 0,
		Duration_s: float64(duration_s),
	}
}

func (ha *HarvestAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	harvestable, _ := ha.Target.(types.IHarvestable); 
	itemHolder, _ := ha.Actor.(types.IItemHolder);
	if itemHolder == nil || harvestable == nil {
		log.Printf("Invalid actor or target in HarvestAction Update()")
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	ha.Timer_s -= dt_s
	if ha.Timer_s <= 0 {
		typeRemoved, amountRemoved := harvestable.Harvest()

		if typeRemoved != "" && amountRemoved > 0 {
			itemHolder.AddItem(typeRemoved, amountRemoved)
			log.Printf("%s added %d %s to inventory", ha.Actor.GetType(), amountRemoved, typeRemoved)
		}

		// harvesting is complete so we return TRUE
		return true
	}

	// harvesting is not complete so we return FALSE
	return false
}