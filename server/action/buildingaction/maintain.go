package buildingaction

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

// "maintain": Restores Pulse to a non-dead building
// - has a set duration and restores a fixed amount of pulse based on gotchi traits
// - consumes a 1 kekwood and 1 alphaslate from actor

type MaintainAction struct {
	action.Action

	PulseRestoredPerSecond int
	LastRestoreTime time.Time

	MaintenanceDuration_s time.Duration
	MaintenanceStartTime time.Time

	TotalPulseRestored int
}

func NewMaintainAction(actor, target types.IEntity, weighting float64) *MaintainAction {
	actorItemHolder, _ := actor.(types.IInventory)
	actorStats, _ := actor.(stats.IStats)
	if actorStats == nil || actorItemHolder == nil {
		log.Printf("ERROR [%s]: Actor does not have IStats or IInventory, returning...", utils.GetFuncName())
		return nil
	}

	// Ecto determines maintenance duration
	actorEcto := actorStats.GetStat(stats.Ecto)
	if actorEcto < 0 {
		log.Printf("ERROR [%s]: Actor does not have ESP stats, returning...", utils.GetFuncName())
		return nil
	}

	// find spark delta from farmer peak and clamp it between 0 and 500
	deltaToPeakEcto := utils.Abs(actorEcto - jobs.Engineer.Peak.Ecto)
	deltaToPeakEcto = utils.Clamp(deltaToPeakEcto, 0, 500)

	// vary pulse restored per sec between 5 and 20
	alpha := float64(deltaToPeakEcto) / 500.0
	pulseRestoredPerSecond := int(5 + 15 * alpha)

	return &MaintainAction{
		Action: action.Action{
			Type: "maintain",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
		PulseRestoredPerSecond: pulseRestoredPerSecond,
		LastRestoreTime: time.Now(),

		MaintenanceDuration_s: time.Duration(30) * time.Second,
		MaintenanceStartTime: time.Now(),
	
		TotalPulseRestored: 0,
	}
}

func (a *MaintainAction) CanBeExecuted() bool {
	maintainable, _ := a.Target.(types.IMaintainable); 
	itemHolder, _ := a.Actor.(types.IInventory);

	// actor and target of correct types?
	if itemHolder == nil || maintainable == nil {
		log.Printf("ERROR [%s]: Invalid actor or target, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(a.Target) {
		return false
	}

	// entity is ready to be maintained?
	if !maintainable.CanBeMaintained() {
		return false
	}

	// actor has 1 kekwood and 1 alphaslate?
	if itemHolder.GetItemQuantity("kekwood") <= 0 || 
	itemHolder.GetItemQuantity("alphaslate") <= 0 {
		return false
	}

	// ok can execute
	return true
}

func (a *MaintainAction) Start() {
	// move to target
	a.TryMoveToTargetEntity(a.Target)
	a.MaintenanceStartTime = time.Now()
}

func (a *MaintainAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	maintainable, _ := a.Target.(types.IMaintainable)
	maintainableStats, _ := a.Target.(stats.IStats)
	itemHolder, _ := a.Actor.(types.IInventory);
	if itemHolder == nil || maintainableStats == nil || maintainable == nil {
		log.Printf("ERROR [%s]: Invalid maintainable, returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	if time.Since(a.LastRestoreTime) > time.Duration(1) * time.Second {
		a.LastRestoreTime = time.Now()

		// do maintenance by adding pulse
		maintainable.Maintain(a.PulseRestoredPerSecond)
		a.TotalPulseRestored += a.PulseRestoredPerSecond

		// check if maintenance is complete due to going over max pulse
		if maintainableStats.GetStat(stats.Pulse) >= maintainable.GetMaxPulse() {
			itemHolder.RemoveItem("kekwood", 1)
			itemHolder.RemoveItem("alphaslate", 1)

			maintainableStats.SetStat(stats.Pulse, maintainable.GetMaxPulse())

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
	if time.Since(a.MaintenanceStartTime) > a.MaintenanceDuration_s {
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