package buildingaction

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

// "maintain"
// 	Restores Pulse to a non-dead building.
//  Action continues until either:
// 		a) building Pulse is fully restored
// 		b) builder runs out of kekwood or alphaslate

type MaintainAction struct {
	action.Action

	PulseRestoredPerSecond float64
	Timer_s float64
	PulseBuffer float64 // this is built up by consuming 1 kekwood and 1 alphaslate

	TotalPulseRestored float64
}

func NewMaintainAction(actor, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *MaintainAction {

	actorItemHolder, _ := actor.(interfaces.IInventory)
	actorStats, _ := actor.(interfaces.IStats)
	if actorStats == nil || actorItemHolder == nil {
		log.Printf("ERROR [%s]: Actor does not have IStats or IInventory, returning...", utils.GetFuncName())
		return nil
	}

	// Ecto determines maintenance duration
	actorEcto := actorStats.GetStat(stattypes.Ecto)
	if actorEcto < 0 {
		log.Printf("ERROR [%s]: Actor does not have ESP stats, returning...", utils.GetFuncName())
		return nil
	}

	// check for gotchi job multiplier
	jobMultiplier, err := utils.GetJobActionMultiplier(actor, "maintain")
	if err != nil {
		log.Printf("ERROR [%s]: Invalid actor or action name, returning...", utils.GetFuncName())
	}

	// vary pulse restored per sec between 0.1 - 1.0
	alpha := actorEcto / 1000
	pulseRestoredPerSecond := (0.1 + 0.9 * alpha) * float64(jobMultiplier)

	wm := actor.GetZone().GetWorldManager()

	a := &MaintainAction{
		Action: action.Action{
			Type: "maintain",
			Weighting: weighting,
			Actor: actor,
			Target: target,
			WorldManager: wm,
		},
		PulseRestoredPerSecond: pulseRestoredPerSecond,
		Timer_s: 0,
		PulseBuffer: 0,
		TotalPulseRestored: 0,
	}

	a.SetFallbackTargetSpec(fallbackTargetSpec)

	return a
}

func (a *MaintainAction) IsValidTarget(potentialTarget interfaces.IEntity) bool {
	if potentialTarget == nil {
		return false
	}

	maintainable, _ := potentialTarget.(interfaces.IMaintainable)
	if maintainable == nil {
		log.Printf("ERROR [%s]: Invalid target, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// entity is ready to be maintained?
	if !maintainable.CanBeMaintained() {
		return false
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(potentialTarget) {
		return false
	}

	return true
}

func (a *MaintainAction) IsValidActor(potentialActor interfaces.IEntity) bool {
	itemHolder, _ := potentialActor.(interfaces.IInventory)
	if itemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// actor has 1 kekwood and 1 alphaslate?
	if itemHolder.GetItemQuantity("kekwood") <= 0 || 
		itemHolder.GetItemQuantity("alphaslate") <= 0 {

		return false
	}

	return true
}

func (a *MaintainAction) Start() {
	// move to target
	a.TryMoveToTargetEntity(a.Target)
}

func (a *MaintainAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	maintainable, _ := a.Target.(interfaces.IMaintainable)
	maintainableStats, _ := a.Target.(interfaces.IStats)
	itemHolder, _ := a.Actor.(interfaces.IInventory)
	actorStats, _ := a.Actor.(interfaces.IStats)
	if itemHolder == nil || maintainableStats == nil || maintainable == nil || actorStats == nil {
		log.Printf("ERROR [%s]: Invalid maintainable, returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// reduce actor ecto and spark
	actorStats.DeltaStat(stattypes.Ecto, -0.1*dt_s)
	actorStats.DeltaStat(stattypes.Spark, -0.1*dt_s)

	// see if 1 second has elapsed
	a.Timer_s -= dt_s
	if a.Timer_s <= 0 {
		a.Timer_s += 1

		// check pulse buffer
		if a.PulseBuffer <= 0 {
			// try top up pulse buffer
			if itemHolder.GetItemQuantity("kekwood") >= 1 && itemHolder.GetItemQuantity("alphaslate") >= 1 {
				a.PulseBuffer += 100
				itemHolder.RemoveItem("kekwood", 1)
				itemHolder.RemoveItem("alphaslate", 1)
			} else {
				// action complete so return true
				a.OnComplete()
				return true
			}
		}

		// do maintenance using PulseRestoredPerSecond
		maintainable.Maintain(a.PulseRestoredPerSecond)
		a.TotalPulseRestored += a.PulseRestoredPerSecond
		a.PulseBuffer -= a.PulseRestoredPerSecond



		// check if maintenance is complete due to going over max pulse
		if maintainableStats.GetStat(stattypes.Pulse) >= 
			maintainableStats.GetStat(stattypes.MaxPulse) {

			// set pulse to max
			maintainableStats.SetStat(stattypes.Pulse, maintainableStats.GetStat(stattypes.MaxPulse))

			// action complete so return true
			a.OnComplete()
			return true
		}
	}

	// maintenance is not complete so we return FALSE
	return false
}

func (a *MaintainAction) OnComplete() {
	itemHolder, _ := a.Actor.(interfaces.IInventory);
	if itemHolder == nil {
		log.Printf("ERROR [%s]: Invalid maintainable, returning...", utils.GetFuncName())
		return	// action is complete we have invalid actor or target
	}

	// 1. grant "buildtoken"s based on TotalPulseRestored
	builderTokenQty := int(a.TotalPulseRestored/100) + 1
	itemHolder.AddItem("buildertoken", builderTokenQty)

	// 2. log activity
	if activityLog, ok := a.Actor.(types.IActivityLog); ok {
		entry := types.ActivityLogEntry{
			Description: fmt.Sprintf("Restored %d Pulse to %s during maintenance and received %d buildertoken's", int(a.TotalPulseRestored), a.Target.GetType(), builderTokenQty),
			LogTime: time.Now(),
		}
		activityLog.NewLogEntry(entry)
	}
}