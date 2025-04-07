package action

import (
	"fmt"
	"log"
	"thereaalm/entity/entitystate"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"
	"thereaalm/utils"
	"time"
)

type AttackAction struct {
	Action
	Timer_s float64
}

func NewAttackAction(actor, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *AttackAction {

	a := &AttackAction{
		Action: Action{
			Type: "attack",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
	}

	a.SetFallbackTargetSpec(fallbackTargetSpec)

	return a
}

func (a *AttackAction) IsValidTarget(potentialTarget interfaces.IEntity) bool {
	if potentialTarget == nil {
		return false
	}

	targetStats, _ := potentialTarget.(interfaces.IStats)
	if targetStats == nil {
		log.Printf("ERROR [%s]: Invalid target, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// can we move to the target?
	if !a.CanMoveToTargetEntity(potentialTarget) {
		return false
	}

	// is target still alive?
	if targetStats.GetStat(stattypes.Pulse) <= 0 {
		return false
	}

	// does target have state and is alive
	targetEntityState, _ := potentialTarget.(entitystate.IEntityState)
	if targetEntityState == nil {
		log.Println("ERROR: An entity without any state was targeted by an attack")
		return false
	}

	if targetEntityState == entitystate.Dead {
		return false
	}

	return true
}

func (a *AttackAction) IsValidActor(potentialActor interfaces.IEntity) bool {
	attackerStats, _ := potentialActor.(interfaces.IStats)

	// do both the target and attacker have stats?
	if attackerStats == nil {
		log.Printf("ERROR [%s]: Invalid actor, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}
	
	// ok we can execute
	return true
}


func (a *AttackAction) Start() {
	// check actor and target are of correct type
	defenderStats, _ := a.Target.(interfaces.IStats)
	attackerStats, _ := a.Actor.(interfaces.IStats)
	if defenderStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction Update()")
		return 	
	}

	// move to target
	a.TryMoveToTargetEntity(a.Target)
}


func (a *AttackAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	defenderStats, _ := a.Target.(interfaces.IStats)
	attackerStats, _ := a.Actor.(interfaces.IStats)
	if defenderStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction Update()")
		return true	// action is complete we have invalid actor or target
	}

	// If target no longer next to us, our action is done
	if !a.Actor.IsNextToTargetEntity(a.Target) {
		return true
	}

	// we attack once per second
	a.Timer_s -= dt_s
	for a.Timer_s <= 0 {
		a.Timer_s += 1

		// attack logic:
		// - spark = attacker power
		// - pulse = defender "HP"

		attackerSpark := attackerStats.GetStat(stattypes.Spark)
		if attackerSpark <= 0 {
			log.Printf("Attacker has no spark to attack with")
			return true
		}

		defenderPulse := defenderStats.GetStat(stattypes.Pulse)
		if defenderPulse <= 0 {
			// log.Printf("Defender has no pulse to defend with")
			return true
		}

		// attacks should be between 1 and 10 attack power for simplicity
		alpha := attackerSpark / 1000
		finalPulseReduction := int(alpha * 10.0)
		finalPulseReduction = utils.Clamp(finalPulseReduction, 1, 10)

		// deal damage to defenders pulse
		defenderStats.DeltaStat(stattypes.Pulse, -finalPulseReduction)
		newDefenderPulse := defenderStats.GetStat(stattypes.Pulse)

		// if defender pulse goes to 0, finish the attack
		if newDefenderPulse <= 0 {
			defenderStats.SetStat(stattypes.Pulse, 0)

			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				entry := types.ActivityLogEntry{
					Description: fmt.Sprintln("Vanquished enemy ", a.Target.GetType()),
					LogTime: time.Now(),
				}
				activityLog.NewLogEntry(entry)
			}

			return true
		}

		// lets make each attack also reduce the attackers ecto by 1 each attack
		attackerStats.DeltaStat(stattypes.Ecto, -1)
		newAttackerEcto := attackerStats.GetStat(stattypes.Ecto)

		if newAttackerEcto < 100 {
			return true
		}

	}

	// attacking is not complete so we return FALSE
	return false
}

func AbsInt(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
