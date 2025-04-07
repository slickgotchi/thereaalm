package interfaces

// IActionPlan is for entities that can process actions.
type IActionPlan interface {
    AddActionToPlan(a IAction)
    SelectNextAction()
    ProcessActions(dt_s float64)
}
