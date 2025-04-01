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

type MineAction struct {
	action.Action

	Duration_s time.Duration
	StartTime time.Time
}

func NewMineAction(actor, target types.IEntity, weighting float64) *MineAction {
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

	// vary duration between 5 - 30 seconds
	alpha := float64(deltaToPeakSpark) / 500.0
	actionDuration_s := int(5 + 25 * alpha)

	return &MineAction{
		Action: action.Action{
			Type: "mine",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
		Duration_s: time.Duration(actionDuration_s) * time.Second,
	}
}

func (a *MineAction) CanBeExecuted() bool {
	mineable, _ := a.Target.(types.IMineable); 
	itemHolder, _ := a.Actor.(types.IInventory);

	// actor and target of correct types?
	if itemHolder == nil || mineable == nil {
		log.Printf("ERROR [%s]: Invalid actor or target, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(a.Target) {
		return false
	}

	// resource entity is ready for collecting?
	if !mineable.CanBeMined() {
		return false
	}

	// ok can execute
	return true
}

func (a *MineAction) Start() {
	// move to target
	if a.TryMoveToTargetEntity(a.Target) {
		a.StartTime = time.Now()
	 }
}

func (a *MineAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	mineable, _ := a.Target.(types.IMineable); 
	itemHolder, _ := a.Actor.(types.IInventory);
	if itemHolder == nil || mineable == nil {
		log.Printf("ERROR [%s]: Invalid actor or target, returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	if time.Since(a.StartTime) > time.Duration(a.Duration_s) {
		typeRemoved, amountRemoved := mineable.Mine()

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
					Description: fmt.Sprintf("Mined %d %s", amountRemoved, typeRemoved),
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