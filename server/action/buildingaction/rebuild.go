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
// - consumes 1 kekwood and 1 alphaslate from actor

type RebuildAction struct {
	action.Action

	PulseRestoredPerSecond int
	LastRestoreTime time.Duration

	RebuildDuration_s time.Duration
	RebuildStartTime time.Duration

	TotalPulseRestored int
}

func NewRebuildAction(actor, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *RebuildAction {
	
	actorItemHolder, _ := actor.(interfaces.IInventory)
	actorStats, _ := actor.(interfaces.IStats)
	if actorStats == nil || actorItemHolder == nil {
		log.Printf("ERROR [%s]: Actor does not have IStats or IInventory, returning...", utils.GetFuncName())
		return nil
	}

	// Ecto determines maintenance duration
	actorPulse := actorStats.GetStat(stattypes.Pulse)
	if actorPulse < 0 {
		log.Printf("ERROR [%s]: Actor does not have ESP stats, returning...", utils.GetFuncName())
		return nil
	}

	// vary pulse restored per sec between 5 and 20
	alpha := actorPulse / 1000
	pulseRestoredPerSecond := int(5 + 15 * alpha)

	wm := actor.GetZone().GetWorldManager()

	a := &RebuildAction{
		Action: action.Action{
			Type: "rebuild",
			Weighting: weighting,
			Actor: actor,
			Target: target,
			WorldManager: wm,
		},
		PulseRestoredPerSecond: pulseRestoredPerSecond,
		LastRestoreTime: 0,

		RebuildDuration_s: time.Duration(30) * time.Second,
		RebuildStartTime: 0,
	
		TotalPulseRestored: 0,
	}

	a.SetFallbackTargetSpec(fallbackTargetSpec)

	return a
}

func (a *RebuildAction) IsValidTarget(potentialTarget interfaces.IEntity) bool {
	if potentialTarget == nil {
		return false
	}

	rebuildable, _ := potentialTarget.(interfaces.IRebuildable); 
	if rebuildable == nil {
		log.Printf("ERROR [%s]: Invalid target, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(potentialTarget) {
		return false
	}

	// entity is ready to be rebuilt?
	if !rebuildable.CanBeRebuilt() {
		return false
	}

	return true
}

func (a *RebuildAction) IsValidActor(potentialActor interfaces.IEntity) bool {
	itemHolder, _ := potentialActor.(interfaces.IInventory);

	// actor and target of correct types?
	if itemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
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
	a.RebuildStartTime = a.WorldManager.Now()
}

func (a *RebuildAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	rebuildable, _ := a.Target.(interfaces.IRebuildable)
	rebuildableStats, _ := a.Target.(interfaces.IStats)
	itemHolder, _ := a.Actor.(interfaces.IInventory);
	if itemHolder == nil || rebuildableStats == nil || rebuildable == nil {
		log.Printf("ERROR [%s]: Invalid actor or target, returning...", utils.GetFuncName())
		return true	// action is complete we have invalid actor or target
	}

	// check duration expired
	if a.WorldManager.Since(a.LastRestoreTime) > 
		time.Duration(1) * time.Second {

		a.LastRestoreTime = a.WorldManager.Now()

		// do rebuild by adding pulse
		rebuildable.Rebuild(a.PulseRestoredPerSecond) 
		a.TotalPulseRestored += a.PulseRestoredPerSecond

		// check if rebuild is complete due to going over max pulse
		if rebuildableStats.GetStat(stattypes.Pulse) >= rebuildableStats.GetStat(stattypes.MaxPulse) {
			itemHolder.RemoveItem("kekwood", 1)
			itemHolder.RemoveItem("alphaslate", 1)

			rebuildableStats.SetStat(stattypes.Pulse, rebuildableStats.GetStat(stattypes.MaxPulse))

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
	if a.WorldManager.Since(a.RebuildStartTime) > 
		a.RebuildDuration_s {

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