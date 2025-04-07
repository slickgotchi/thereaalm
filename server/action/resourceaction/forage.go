package resourceaction

import (
	"fmt"
	"log"
	"thereaalm/action"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"
	"thereaalm/utils"
	"time"
)

type ForageAction struct {
	action.Action

	Duration_s time.Duration
	StartTime time.Time
}

func NewForageAction(actor, target interfaces.IEntity, weighting float64,
		fallbackTargetSpec *types.TargetSpec) *ForageAction {

	actorItemHolder, _ := actor.(types.IInventory)
	actorStats, _ := actor.(interfaces.IStats)
	if actorStats == nil || actorItemHolder == nil {
		log.Printf("ERROR [%s]: Actor does not have IStats or IInventory, returning...", utils.GetFuncName())
		return nil
	}

	// Spark determines gather duration
	actorSpark := actorStats.GetStat(stattypes.Spark)
	if actorSpark < 0 {
		log.Printf("ERROR [%s]: Actor does not have ESP stats, returning...", utils.GetFuncName())
		return nil
	}

	// vary forage duration between 5 - 30 seconds
	alpha := actorSpark / 1000
	actionDuration_s := int(5 + 25 * alpha)

	a := &ForageAction{
		Action: action.Action{
			Type: "forage",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
		Duration_s: time.Duration(actionDuration_s) * time.Second,
	}

	a.SetFallbackTargetSpec(fallbackTargetSpec)

    return a
}

func (a *ForageAction) IsValidTarget(potentialTarget interfaces.IEntity) bool {
	if potentialTarget == nil {
		return false
	}

	forageable, _ := potentialTarget.(types.IForageable); 
	if forageable == nil {
		log.Printf("ERROR [%s]: Invalid target, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(potentialTarget) {
		return false
	}

	// resource entity is ready for collecting?
	if !forageable.CanBeForaged() {
		return false
	}

	return true
}

func (a *ForageAction) IsValidActor(potentialActor interfaces.IEntity) bool {
	itemHolder, _ := potentialActor.(types.IInventory);

	// actor and target of correct types?
	if itemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// ok can execute
	return true
}

func (a *ForageAction) Start() {
	// move to target
	if a.TryMoveToTargetEntity(a.Target) {
		a.StartTime = time.Now()
	 }
}

func (a *ForageAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	forageable, _ := a.Target.(types.IForageable); 
	itemHolder, _ := a.Actor.(types.IInventory);
	if itemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor (not IInventory), returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}
	if forageable == nil {
		log.Printf("ERROR [%s]: Invalid target (not IForageable), returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	if time.Since(a.StartTime) > time.Duration(a.Duration_s) {
		typeRemoved, amountRemoved := forageable.Forage()

		if typeRemoved != "" && amountRemoved > 0 {
			// add item to item holder
			itemHolder.AddItem(typeRemoved, amountRemoved)
	
			// remove some spark and pulse
			if actorStats, ok := a.Actor.(interfaces.IStats); ok {
				actorStats.DeltaStat(stattypes.Spark, -1)
				actorStats.DeltaStat(stattypes.Pulse, -1)
			}
	
			// see if actor has an activity log
			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				entry := types.ActivityLogEntry{
					Description: fmt.Sprintf("Foraged %d %s", amountRemoved, typeRemoved),
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