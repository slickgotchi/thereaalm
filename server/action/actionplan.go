package action

import (
	"math/rand"
	"thereaalm/action/actiontargeting"
	"thereaalm/interfaces"
)

type ActionPlan struct {
    Actions []interfaces.IAction
    CurrentAction interfaces.IAction
}

func (a *ActionPlan) AddActionToPlan(action interfaces.IAction) {
    a.Actions = append(a.Actions, action)
}

func (a *ActionPlan) ProcessActions(dt_s float64) {
	if a.CurrentAction == nil {
        a.SelectNextAction()
        return // Early return if no action to process
    }
    actor := a.CurrentAction.GetActor()
    scaledDt := dt_s
    if consumer, ok := actor.(interfaces.IBuffConsumer); ok {
        scaledDt = dt_s * consumer.GetEffectiveSpeedMultiplier()
    }
    actionComplete := a.CurrentAction.Update(scaledDt)
    if actionComplete {
        a.CurrentAction = nil
    }
	/*
	// If there's no current action, select one based on weightings.
	if a.CurrentAction == nil {
		a.SelectNextAction()
	}

	if a.CurrentAction != nil {
		// Call Update() on the current action to progress it.
		actionComplete := a.CurrentAction.Update(dt_s)

		// If the action is complete, clear it and select a new one.
		if actionComplete {
			// log.Println(a.CurrentAction.GetType(), " complete.")
			a.CurrentAction = nil
		}
	}
		*/
}

// SelectNextAction will only select actions that can be executed.
func (a *ActionPlan) SelectNextAction() {
	// log.Println("Select next action...")

	// If no actions exist, just return.
	if len(a.Actions) == 0 {
		return
	}

	// Calculate total weighting to normalize probabilities, considering only executable actions.
	var totalWeight float64
	executableActions := []interfaces.IAction{}

	// Filter out actions that cannot be executed.
	for _, action := range a.Actions {
		actionCurrentTarget := action.GetTarget()

		// first establish fallbacks if required/possible
		if 	actionCurrentTarget == nil || 
			!action.IsValidTarget(actionCurrentTarget) {
			
				newTarget := actiontargeting.ResolveFallbackTarget(action)
			action.SetTarget(newTarget)
		}

		// see if actor and target is valid
		if !action.IsValidActor(action.GetActor()) || 
			!action.IsValidTarget(action.GetTarget()) {

			continue
		}

		// add to possible executable actions (if possible)
		if action.GetWeighting() > 0 {
			executableActions = append(executableActions, action)
			totalWeight += action.GetWeighting()
		}
	}

	// If no executable actions, return.
	if len(executableActions) == 0 {
		// log.Println("No executable actions available.")
		return
	}

	// Choose a random executable action based on weighting.
	randomWeight := rand.Float64() * totalWeight
	var cumulativeWeight float64
	for _, action := range executableActions {
		cumulativeWeight += action.GetWeighting()
		if cumulativeWeight >= randomWeight {
			a.CurrentAction = action
			a.CurrentAction.Start()
			return
		}
	}
}



// New reporting struct
type ActionPlanReporting struct {
	Actions       []ActionReporting `json:"actions"`
	CurrentAction *ActionReporting  `json:"currentAction,omitempty"`
}

type ActionReporting struct {
	Type       string `json:"type"`
	ActorType  string `json:"actorType"`
	ActorID    string `json:"actorId"`
	TargetType string `json:"targetType,omitempty"`
	TargetID   string `json:"targetId,omitempty"`
	Weighting float64 `json:"weighting"`
}

// ToReporting converts ActionPlan to a cycle-free reporting version
func (a *ActionPlan) ToReporting() ActionPlanReporting {
	actions := make([]ActionReporting, len(a.Actions))
	for i, action := range a.Actions {
		var targetType, targetID string
		if action.GetTarget() != nil {
			targetType = action.GetTarget().GetType()
			targetID = action.GetTarget().GetUUID().String()
		}
		actions[i] = ActionReporting{
			Type:       action.GetType(),
			ActorType:  action.GetActor().GetType(),
			ActorID:    action.GetActor().GetUUID().String(),
			TargetType: targetType,
			TargetID:   targetID,
			Weighting: action.GetWeighting(),
		}
	}

	var current *ActionReporting
	if a.CurrentAction != nil {
		var targetType, targetID string
		if a.CurrentAction.GetTarget() != nil {
			targetType = a.CurrentAction.GetTarget().GetType()
			targetID = a.CurrentAction.GetTarget().GetUUID().String()
		}
		current = &ActionReporting{
			Type:       a.CurrentAction.GetType(),
			ActorType:  a.CurrentAction.GetActor().GetType(),
			ActorID:    a.CurrentAction.GetActor().GetUUID().String(),
			TargetType: targetType,
			TargetID:   targetID,
			Weighting: a.CurrentAction.GetWeighting(),
		}
	}

	return ActionPlanReporting{
		Actions:       actions,
		CurrentAction: current,
	}
}