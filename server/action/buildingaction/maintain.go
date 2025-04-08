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

// "maintain": Restores Pulse to a non-dead building
// - has a set duration and restores a fixed amount of pulse based on gotchi traits
// - consumes a 1 kekwood and 1 alphaslate from actor

type MaintainAction struct {
	action.Action

	PulseRestoredPerSecond int
	LastRestoreTime time.Duration

	MaintenanceDuration_s time.Duration
	MaintenanceStartTime time.Duration

	TotalPulseRestored int
}

func NewMaintainAction(actor, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *MaintainAction {

	actorItemHolder, _ := actor.(types.IInventory)
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

	// vary pulse restored per sec between 5 and 20
	alpha := actorEcto / 1000
	pulseRestoredPerSecond := int(5 + 15 * alpha)

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
		LastRestoreTime: 0,

		MaintenanceDuration_s: time.Duration(30) * time.Second,
		MaintenanceStartTime: 0,
	
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
	itemHolder, _ := potentialActor.(types.IInventory)
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
	a.MaintenanceStartTime = a.WorldManager.Now()
}

func (a *MaintainAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	maintainable, _ := a.Target.(interfaces.IMaintainable)
	maintainableStats, _ := a.Target.(interfaces.IStats)
	itemHolder, _ := a.Actor.(types.IInventory);
	if itemHolder == nil || maintainableStats == nil || maintainable == nil {
		log.Printf("ERROR [%s]: Invalid maintainable, returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	if a.WorldManager.Since(a.LastRestoreTime) > 
		time.Duration(1) * time.Second {

		a.LastRestoreTime = a.WorldManager.Now()

		// do maintenance by adding pulse
		maintainable.Maintain(a.PulseRestoredPerSecond)
		a.TotalPulseRestored += a.PulseRestoredPerSecond

		// check if maintenance is complete due to going over max pulse
		if maintainableStats.GetStat(stattypes.Pulse) >= maintainable.GetMaxPulse() {
			itemHolder.RemoveItem("kekwood", 1)
			itemHolder.RemoveItem("alphaslate", 1)

			maintainableStats.SetStat(stattypes.Pulse, maintainable.GetMaxPulse())

			// see if actor has an activity log
			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				entry := types.ActivityLogEntry{
					Description: fmt.Sprintf("Restored %d Pulse to %s during maintenance", a.TotalPulseRestored, a.Target.GetType()),
					LogTime: time.Now(),
				}
				activityLog.NewLogEntry(entry)
			}

			// maintenance is complete so return true
			return true
		}
	}

	// check maintenance duration expired?
	if a.WorldManager.Since(a.MaintenanceStartTime) > 
		a.MaintenanceDuration_s {

		itemHolder.RemoveItem("kekwood", 1)
		itemHolder.RemoveItem("alphaslate", 1)

		// see if actor has an activity log
		if activityLog, ok := a.Actor.(types.IActivityLog); ok {
			entry := types.ActivityLogEntry{
				Description: fmt.Sprintf("Restored %d Pulse to %s during maintenance", a.TotalPulseRestored, a.Target.GetType()),
				LogTime: time.Now(),
			}
			activityLog.NewLogEntry(entry)
		}

		// maintenance is complete so return true
		return true
	}

	// harvesting is not complete so we return FALSE
	return false
}