package resourceaction

import (
	"fmt"
	"log"
	"thereaalm/action"
	"thereaalm/jobs"
	"thereaalm/stats"
	"thereaalm/types"
	"thereaalm/utils"
	"time"
)

type ForageAction struct {
	action.Action

	Duration_s time.Duration
	StartTime time.Time
}

func NewForageAction(actor, target types.IEntity, weighting float64,
		fallbackTargetSpec *action.TargetSpec) *ForageAction {

	actorItemHolder, _ := actor.(types.IInventory)
	actorStats, _ := actor.(stats.IStats)
	if actorStats == nil || actorItemHolder == nil {
		log.Printf("ERROR [%s]: Actor does not have IStats or IInventory, returning...", utils.GetFuncName())
		return nil
	}

	// Spark determines gather duration
	actorSpark := actorStats.GetStat(stats.Spark)
	if actorSpark < 0 {
		log.Printf("ERROR [%s]: Actor does not have ESP stats, returning...", utils.GetFuncName())
		return nil
	}

	// find spark delta from farmer peak and clamp it between 0 and 500
	deltaToPeakSpark := utils.Abs(actorSpark - jobs.Farmer.Peak.Spark)
	deltaToPeakSpark = utils.Clamp(deltaToPeakSpark, 0, 500)

	// vary forage duration between 5 - 30 seconds
	alpha := float64(deltaToPeakSpark) / 500.0
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

func (a *ForageAction) CanBeExecuted() bool {
	// check current target validity and/or set a fallback if neccessary
	if !a.EnsureValidTarget() {
		return false
	}

	forageable, _ := a.Target.(types.IForageable); 
	itemHolder, _ := a.Actor.(types.IInventory);

	// actor and target of correct types?
	if itemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor (not IInventory), returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}
	if forageable == nil {
		log.Printf("ERROR [%s]: Invalid target (not IForageable), returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(a.Target) {
		return false
	}

	// resource entity is ready for collecting?
	if !forageable.CanBeForaged() {
		return false
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
			if actorStats, ok := a.Actor.(stats.IStats); ok {
				actorStats.DeltaStat(stats.Spark, -1)
				actorStats.DeltaStat(stats.Pulse, -1)
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