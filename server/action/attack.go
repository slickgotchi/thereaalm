package action

import (
	"fmt"
	"log"
	"thereaalm/jobs"
	"thereaalm/mathext"
	"thereaalm/stats"
	"thereaalm/types"
	"time"
)

type AttackAction struct {
	Action
	Timer_s float64
}

func NewAttackAction(actor, target types.IEntity, weighting float64) *AttackAction {
	return &AttackAction{
		Action: Action{
			Type: "attack",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
	}
}

func (a *AttackAction) CanBeExecuted() bool {
	targetStats, _ := a.Target.(stats.IStats)
	attackerStats, _ := a.Actor.(stats.IStats)

	// do both the target and attacker have stats?
	if targetStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction CanBeExecuted()")
		return false // can not execute, invalid actor or target
	}

	// can we move to the target?
	if !a.CanMoveToTargetEntity(a.Target) {
		return false
	}

	// is target still alive?
	if targetStats.GetStat(stats.Spark) <= 0 {
		return false
	}
	
	// ok we can execute
	return true
}


func (a *AttackAction) Start() {
	// check actor and target are of correct type
	defenderStats, _ := a.Target.(stats.IStats)
	attackerStats, _ := a.Actor.(stats.IStats)
	if defenderStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction Update()")
		return 	
	}

	// move to target
	a.TryMoveToTargetEntity(a.Target)
}


func (a *AttackAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	defenderStats, _ := a.Target.(stats.IStats)
	attackerStats, _ := a.Actor.(stats.IStats)
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
		// - pulse => attack power => reduces defenders spark
		// - ecto => crit chance => improves attack power (sometimes) 

		attackerPulse := attackerStats.GetStat(stats.Pulse)
		if attackerPulse <= 0 {
			log.Printf("Attacker has no pulse to attack with")
			return true
		}

		defemderSpark := defenderStats.GetStat(stats.Spark)
		if defemderSpark <= 0 {
			log.Printf("Defender has no spark to defend with")
			return true
		}

		// use mercenary peak pulse to determine attack power
		deltaPulse := mathext.Abs(attackerPulse - jobs.Mercenary.Peak.Pulse)
		
		// clamp between 0 and 500
		deltaPulse = mathext.Clamp(deltaPulse, 0, 500)

		// attacks should be between 1 and 10 attack power for simplicity
		alpha := float64(500 - deltaPulse) / 500
		finalSparkReduction := int(alpha * 10.0)
		finalSparkReduction = mathext.Clamp(finalSparkReduction, 1, 10)

		// deal damage to defenders spark
		defenderStats.DeltaStat(stats.Spark, -finalSparkReduction)
		newDefenderSpark := defenderStats.GetStat(stats.Spark)

		// ecto goes below 10% (100) we need to finish the attack
		if newDefenderSpark <= 0 {
			defenderStats.SetStat(stats.Spark, 0)
			log.Println("Defeated enemy")

			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				entry := types.ActivityLogEntry{
					Description: fmt.Sprintf("Vanquished enemy ", a.Target.GetType()),
					LogTime: time.Now(),
				}
				activityLog.NewLogEntry(entry)
			}

			return true
		}

		// lets make each attack also reduce the attackers ecto by 1 each attack
		attackerStats.DeltaStat(stats.Ecto, -1)
		newAttackerEcto := attackerStats.GetStat(stats.Ecto)

		if newAttackerEcto < 100 {
			return true
		}

	}

	// harvesting is not complete so we return FALSE
	return false
}

func AbsInt(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
