package types

// IAction defines an executable behavior.
type IAction interface {
    Start()
    CanBeExecuted() bool
    Update(dt_s float64) bool

    // these are defined in Action (should not need to override)
    GetType() string
    GetWeighting() float64
    GetTarget() IEntity
    GetActor() IEntity
}

// IActionPlan is for entities that can process actions.
type IActionPlan interface {
    AddActionToPlan(a IAction)
    SelectNextAction()
    ProcessActions(dt_s float64)
}

