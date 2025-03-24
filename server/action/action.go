package action

import (
	"log"
	"math/rand"
	"thereaalm/types"
)

type Action struct {
	Type string
	Weighting float64
	// IsStarted bool
	Actor types.IEntity
	Target types.IEntity
}

// this should be overidden
// return TRUE when the action is COMPLETE
func (a *Action) Update(dt_s float64) bool { return true }

// don't need to be overridden
func (a *Action) GetType() string {return a.Type}
func (a *Action) GetWeighting() float64 {return a.Weighting}
func (a *Action) GetTarget() types.IEntity {return a.Target}
func (a *Action) GetActor() types.IEntity {return a.Actor}
func (a *Action) CanBeExecuted() bool {return true}

type ActionPlan struct {
    Actions []types.IAction
    CurrentAction types.IAction
}

func (a *ActionPlan) AddAction(action types.IAction) {
    a.Actions = append(a.Actions, action)
}

// SelectNextAction will only select actions that can be executed.
func (a *ActionPlan) SelectNextAction() {
	log.Println("Select next action...")

	// If no actions exist, just return.
	if len(a.Actions) == 0 {
		return
	}

	// Calculate total weighting to normalize probabilities, considering only executable actions.
	var totalWeight float64
	executableActions := []types.IAction{}

	// Filter out actions that cannot be executed.
	for _, action := range a.Actions {
		if action.CanBeExecuted() && action.GetWeighting() > 0 {
			executableActions = append(executableActions, action)
			totalWeight += action.GetWeighting()
		}
	}

	// If no executable actions, return.
	if len(executableActions) == 0 {
		log.Println("No executable actions available.")
		return
	}

	// Choose a random executable action based on weighting.
	randomWeight := rand.Float64() * totalWeight
	var cumulativeWeight float64
	for _, action := range executableActions {
		cumulativeWeight += action.GetWeighting()
		if cumulativeWeight >= randomWeight {
			a.CurrentAction = action
			log.Println("Selected action: ", action.GetType())
			return
		}
	}
}

func (a *ActionPlan) ProcessActions(dt_s float64) {
	// If there's no current action, select one based on weightings.
	if a.CurrentAction == nil {
		a.SelectNextAction()
	}

	if a.CurrentAction != nil {
		// Call Update() on the current action to progress it.
		actionComplete := a.CurrentAction.Update(dt_s)

		// If the action is complete, clear it and select a new one.
		if actionComplete {
			log.Println(a.CurrentAction.GetType(), " complete.")
			a.CurrentAction = nil
		}
	}
}