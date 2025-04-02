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

// utility function to move to a target
func (a *Action) CanMoveToTargetEntity(target types.IEntity) bool {
	zone := a.Actor.GetZone()
	_, _, found := zone.TryGetEmptyTileNextToTargetEntity(target)
	return found
}

func (a *Action) TryMoveToTargetEntity(target types.IEntity) bool {
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
			a.CurrentAction.Start()
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