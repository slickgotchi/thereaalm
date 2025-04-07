package action

import (
	"log"
	// "thereaalm/action/actiontargeting"
	"thereaalm/interfaces"
	"thereaalm/types"
)

type Action struct {
	Type string
	Weighting float64
	Actor interfaces.IEntity
	Target interfaces.IEntity
	FallbackTargetSpec *types.TargetSpec // Fallback target specification
}

func (a *Action) Start() {}
func (a *Action) Update(dt_s float64) bool { return true }

func (a *Action) IsValidTarget(potentialTarget interfaces.IEntity) bool { 
	log.Println("WARNING: IsValidTarget has not been overridden by action ", a.Type)
	return potentialTarget != nil 
}
func (a *Action) IsValidActor(potentialActor interfaces.IEntity) bool { 
	log.Println("WARNING: IsValidActor has not been overridden by action ", a.Type)
	return potentialActor != nil 
}

// func (a *Action) TryUseFallbackIfTargetIsNil() bool {
// 	// Validate current target
// 	if a.Target != nil {
// 		return true
// 	}

// 	// Attempt fallback resolution if no valid target
// 	if a.FallbackFunc != nil && a.FallbackTargetSpec != nil {
// 		if newTarget := a.FallbackFunc(a); newTarget != nil {
// 			a.Target = newTarget
// 			return true
// 		}
// 	}
	
// 	// could not even use a fallback
// 	log.Printf("Fallback target resolution failed for action %s", a.Type)
// 	return false
// }

func (a *Action) GetType() string {return a.Type}
func (a *Action) GetWeighting() float64 {return a.Weighting}
func (a *Action) GetTarget() interfaces.IEntity {return a.Target}
func (a *Action) GetActor() interfaces.IEntity {return a.Actor}
func (a *Action) SetTarget(newTarget interfaces.IEntity) {
	a.Target = newTarget
}

func (a *Action) GetFallbackTargetSpec() *types.TargetSpec {
	return a.FallbackTargetSpec
}

// SetFallbackTargetSpec sets the fallback target specification and corresponding fallback function
func (a *Action) SetFallbackTargetSpec(fallbackTargetSpec *types.TargetSpec) {
    if fallbackTargetSpec == nil {
        return
    }
    a.FallbackTargetSpec = fallbackTargetSpec
    // if handler, exists := actiontargeting.FallbackHandlers[fallbackTargetSpec.TargetCriterion]; exists {
    //     a.FallbackFunc = handler
    // } else {
    //     log.Printf("No fallback handler for criterion %s in action %s", fallbackTargetSpec.TargetCriterion, a.Type)
    // }
}

// utility function to move to a target
func (a *Action) CanMoveToTargetEntity(target interfaces.IEntity) bool {
	zone := a.Actor.GetZone()
	_, _, found := zone.TryGetEmptyTileNextToTargetEntity(target)
	return found
}

func (a *Action) TryMoveToTargetEntity(target interfaces.IEntity) bool {
	// check if already next to target
	if a.Actor.IsNextToTargetEntity(target) {
		// ensure we're facing the target
		a.Actor.SetDirectionToTargetEntity(target)
		return true
	} else {
		// check spatial map for a valid position next to the target
		zone := a.Actor.GetZone()
		nx, ny, found := zone.TryGetEmptyTileNextToTargetEntity(target)
		if !found {
			return false
		}

		// move to target
		a.Actor.SetPosition(nx, ny)

		// ensure we're facing the target
		a.Actor.SetDirectionToTargetEntity(target)

		return true
	}
}

func (a *Action) CanMoveToTargetPosition(x, y int) bool {
	zone := a.Actor.GetZone()
	return !zone.IsTileOccupied(x,y)
}

func (a *Action) TryMoveToTargetPosition(x, y int) bool {
	currX, currY := a.Actor.GetPosition()

	// check if already at the target position
	if currX == x && currY == y {
		return true
	} else {
		// check spatial map that target is valid position
		zone := a.Actor.GetZone()
		if zone.IsTileOccupied(x, y) {
			return false
		}

		// ensure we're facing the target position
		a.Actor.SetDirectionToTargetPosition(x, y)

		// move to target
		a.Actor.SetPosition(x, y)

		return true
	}
}