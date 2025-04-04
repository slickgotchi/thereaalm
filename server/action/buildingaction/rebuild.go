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

type RebuildAction struct {
	action.Action

	PulseRestoredPerSecond int
	LastRestoreTime time.Time

	RebuildDuration_s time.Duration
	RebuildStartTime time.Time

	TotalPulseRestored int
}

func NewRebuildAction(actor, target types.IEntity, weighting float64) *RebuildAction {
	actorItemHolder, _ := actor.(types.IInventory)
	actorStats, _ := actor.(stats.IStats)
	if actorStats == nil || actorItemHolder == nil {
		log.Printf("ERROR [%s]: Actor does not have IStats or IInventory, returning...", utils.GetFuncName())
		return nil
	}

	// Ecto determines maintenance duration
	actorPulse := actorStats.GetStat(stats.Pulse)
	if actorPulse < 0 {
		log.Printf("ERROR [%s]: Actor does not have ESP stats, returning...", utils.GetFuncName())
		return nil
	}

	// find spark delta from farmer peak and clamp it between 0 and 500
	deltaToPeakPulse := utils.Abs(actorPulse - jobs.Builder.Peak.Pulse)
	deltaToPeakPulse = utils.Clamp(deltaToPeakPulse, 0, 500)

	// vary pulse restored per sec between 5 and 20
	alpha := float64(deltaToPeakPulse) / 500.0
	pulseRestoredPerSecond := int(5 + 15 * alpha)

	return &RebuildAction{
		Action: action.Action{
			Type: "rebuild",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
		PulseRestoredPerSecond: pulseRestoredPerSecond,
		LastRestoreTime: time.Now(),

		RebuildDuration_s: time.Duration(30) * time.Second,
		RebuildStartTime: time.Now(),
	
		TotalPulseRestored: 0,
	}
}

func (a *RebuildAction) CanBeExecuted() bool {
	rebuildable, _ := a.Target.(types.IRebuildable); 
	itemHolder, _ := a.Actor.(types.IInventory);

	// actor and target of correct types?
	if itemHolder == nil || rebuildable == nil {
		log.Printf("ERROR [%s]: Invalid actor or target, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(a.Target) {
		return false
	}

	// entity is ready to be maintained?
	if !rebuildable.CanBeRebuilt() {
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

func (a *RebuildAction) Start() {
	// move to target
	a.TryMoveToTargetEntity(a.Target)
	a.RebuildStartTime = time.Now()
}

func (a *RebuildAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	rebuildable, _ := a.Target.(types.IRebuildable)
	rebuildableStats, _ := a.Target.(stats.IStats)
	itemHolder, _ := a.Actor.(types.IInventory);
	if itemHolder == nil || rebuildableStats == nil || rebuildable == nil {
		log.Printf("ERROR [%s]: Invalid actor or target, returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	if time.Since(a.LastRestoreTime) > time.Duration(1) * time.Second {
		a.LastRestoreTime = time.Now()

		// do rebuild by adding pulse
		rebuildable.Rebuild(a.PulseRestoredPerSecond)
		a.TotalPulseRestored += a.PulseRestoredPerSecond

		// check if rebuild is complete due to going over max pulse
		if rebuildableStats.GetStat(stats.Pulse) >= rebuildable.GetMaxPulse() {
			itemHolder.RemoveItem("kekwood", 1)
			itemHolder.RemoveItem("alphaslate", 1)

			rebuildableStats.SetStat(stats.Pulse, rebuildable.GetMaxPulse())

			// see if actor has an activity log
			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				entry := types.ActivityLogEntry{
					Description: fmt.Sprintf("Restored %d Pulse to %s during rebuild", a.TotalPulseRestored, a.Target.GetType()),
					LogTime: time.Now(),
				}
				activityLog.NewLogEntry(entry)
			}

			// maintenance is complete so return true
			return true
		}
	}

	// check maintenance duration expired?
	if time.Since(a.RebuildStartTime) > a.RebuildDuration_s {
		itemHolder.RemoveItem("kekwood", 1)
		itemHolder.RemoveItem("alphaslate", 1)

		// see if actor has an activity log
		if activityLog, ok := a.Actor.(types.IActivityLog); ok {
			entry := types.ActivityLogEntry{
				Description: fmt.Sprintf("Restored %d Pulse to %s during rebuild", a.TotalPulseRestored, a.Target.GetType()),
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