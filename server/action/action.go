package action

import (
	"thereaalm/types"
)

type Action struct {
	Type string
	IsStarted bool
	Actor types.IEntity
	Target types.IEntity
}

func (a *Action) Execute(actor, target types.IEntity) bool {
	a.IsStarted = true
	a.Actor = actor
	a.Target = target
	return true
}

// this should be overidden
// return TRUE when the action is COMPLETE
func (a *Action) Update(dt_s float64) bool { return true }

// don't need to be overridden
func (a *Action) GetType() string {return a.Type}
func (a *Action) GetTarget() types.IEntity {return a.Target}
func (a *Action) GetActor() types.IEntity {return a.Actor}

