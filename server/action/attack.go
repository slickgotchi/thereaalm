package action

import (
	"fmt"
	"log"
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
	if targetStats.GetStat(stats.HpCurrent) <= 0 {
		return false
	}
	
	// ok we can execute
	return true
}


func (a *AttackAction) Start() {
	// check actor and target are of correct type
	targetStats, _ := a.Target.(stats.IStats)
	attackerStats, _ := a.Actor.(stats.IStats)
	if targetStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction Update()")
		return 	
	}

	// move to target
	a.TryMoveToTargetEntity(a.Target)
}


func (a *AttackAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	targetStats, _ := a.Target.(stats.IStats)
	attackerStats, _ := a.Actor.(stats.IStats)
	if targetStats == nil || attackerStats == nil {
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

		attackDamage:= attackerStats.GetStat(stats.Attack); 
		if attackDamage <= 0 {
			log.Printf("attacker does not have attack stat")
			return true
		}

		targetStats.DeltaStat(stats.HpCurrent, -attackDamage)
		newHp := targetStats.GetStat(stats.HpCurrent)

		if newHp <= 0 {
			targetStats.SetStat(stats.HpCurrent, 0)
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
	}

	// harvesting is not complete so we return FALSE
	return false
}