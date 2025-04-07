package interfaces

import (
	"thereaalm/types"
)

// IAction defines an executable behavior.
type IAction interface {
    Start()
    Update(dt_s float64) bool

    IsValidTarget(potentialTarget IEntity) bool
    IsValidActor(potentialActor IEntity) bool

    // TryUseFallbackIfTargetIsNil() bool

    GetType() string
    GetWeighting() float64
    GetTarget() IEntity
    GetActor() IEntity
    SetTarget(newTarget IEntity)

    GetFallbackTargetSpec() *types.TargetSpec
    SetFallbackTargetSpec(fallbackTargetspec *types.TargetSpec)

    CanMoveToTargetEntity(target IEntity) bool
    TryMoveToTargetEntity(target IEntity) bool
    CanMoveToTargetPosition(x, y int) bool 
    TryMoveToTargetPosition(x, y int) bool 
}