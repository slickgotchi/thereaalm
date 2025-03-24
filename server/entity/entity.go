package entity

import (
	"thereaalm/types"

	"github.com/google/uuid"
)

// ENTITY
type Entity struct {
	ID uuid.UUID
	Type string
    X int
    Y int
}

func (e *Entity) GetUUID() uuid.UUID { return e.ID }
func (e *Entity) GetType() string          { return e.Type }
func (e *Entity) Update(dt_s float64) {}
func (e *Entity) GetPosition() (int, int) {
	return e.X, e.Y
}
func (e *Entity) SetPosition(x, y int) {
	e.X = x
	e.Y = y
}
func (e *Entity) GetCustomData() interface{} {
    return nil
}


type ActionSequence struct {
    Actions []types.IAction
}

func (a *ActionSequence) QueueAction(action types.IAction) {
    a.Actions = append(a.Actions, action)
}

func (a *ActionSequence) ProcessActions(dt_s float64) {
    for _, action := range a.Actions {
        action.Execute(action.GetActor(), action.GetTarget()) // You can replace `nil` with actual actor/target when needed
    }
    a.Actions = nil // Clear actions after processing, if needed
}