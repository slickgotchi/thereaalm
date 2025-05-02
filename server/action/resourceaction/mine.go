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

type MineAction struct {
	action.Action

	Timer_s float64
}

func NewMineAction(actor, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *MineAction {

	actorItemHolder, _ := actor.(interfaces.IInventory)
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

	// check for gotchi job multiplier
	jobMultiplier, err := utils.GetJobActionMultiplier(actor, "mine")
	if err != nil {
		log.Printf("ERROR [%s]: Invalid actor or action name, returning...", utils.GetFuncName())
	}

	// vary duration between 150 - 300 seconds
	alpha := actorSpark / 1000
	actionDuration_s := (150 + 150 * (1-alpha)) / jobMultiplier

	wm := actor.GetZone().GetWorldManager()

	a := &MineAction{
		Action: action.Action{
			Type: "mine",
			Weighting: weighting,
			Actor: actor,
			Target: target,
			WorldManager: wm,
		},
		Timer_s: float64(actionDuration_s),
	}

	a.SetFallbackTargetSpec(fallbackTargetSpec)

	return a
}

func (a *MineAction) IsValidTarget(potentialTarget interfaces.IEntity) bool {
	if potentialTarget == nil {
		return false
	}

	mineable, _ := potentialTarget.(types.IMineable)
	if mineable == nil {
		log.Printf("ERROR [%s]: Invalid target, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(potentialTarget) {
		return false
	}

	// resource entity is ready for collecting?
	if !mineable.CanBeMined() {
		return false
	}

	return true
}

func (a *MineAction) IsValidActor(potentialActor interfaces.IEntity) bool {
	itemHolder, _ := potentialActor.(interfaces.IInventory);

	// actor and target of correct types?
	if itemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// ok can execute
	return true
}

func (a *MineAction) Start() {
	// move to target
	a.TryMoveToTargetEntity(a.Target)
}

func (a *MineAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	mineable, _ := a.Target.(types.IMineable); 
	itemHolder, _ := a.Actor.(interfaces.IInventory);
	if itemHolder == nil || mineable == nil {
		log.Printf("ERROR [%s]: Invalid actor or target, returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// remove some spark and pulse
	if actorStats, ok := a.Actor.(interfaces.IStats); ok {
		actorStats.DeltaStat(stattypes.Spark, -0.1*dt_s)
		actorStats.DeltaStat(stattypes.Pulse, -0.1*dt_s)
	}

	// check duration expired
	a.Timer_s -= dt_s
	if a.Timer_s <= 0 {

		typeRemoved, amountRemoved := mineable.Mine()

		if typeRemoved != "" && amountRemoved > 0 {
			// add item to item holder
			itemHolder.AddItem(typeRemoved, amountRemoved)
	
			// see if actor has an activity log
			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				entry := types.ActivityLogEntry{
					Description: fmt.Sprintf("Mined %d %s", amountRemoved, typeRemoved),
					LogTime: time.Now(),
				}
				activityLog.NewLogEntry(entry)
			}
		}
	
		// action is complete so we return TRUE
		return true
	}

	// action is not complete so we return FALSE
	return false
}