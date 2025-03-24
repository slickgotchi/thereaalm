package types

// IAction defines an executable behavior.
type IAction interface {
    // Execute(actor, target IEntity) bool
    Update(dt_s float64) bool
    GetType() string
    GetWeighting() float64
    GetTarget() IEntity
    GetActor() IEntity
    CanBeExecuted() bool
}

// IActionPlan is for entities that can process actions.
type IActionPlan interface {
    AddAction(a IAction)
    SelectNextAction()
    ProcessActions(dt_s float64)
}