package types

// IAction defines an executable behavior.
type IAction interface {
    Execute(actor, target IEntity) bool
    Update(dt_s float64) bool
    GetType() string
    GetTarget() IEntity
    GetActor() IEntity
}