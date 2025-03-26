package action

import (
	"log"
	"thereaalm/stats"
	"thereaalm/types"
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

func NextToEachOther(a, b types.IEntity) bool {
    ax, ay := a.GetPosition()
    bx, by := b.GetPosition()

    // Check if the entities are next to each other (left, right, up, down)
    return (ax == bx && (ay == by+1 || ay == by-1)) || // Vertical check (up, down)
           (ay == by && (ax == bx+1 || ax == bx-1))   // Horizontal check (left, right)
}

func (action *AttackAction) Start() {
	// check actor and target are of correct type
	targetStats, _ := action.Target.(stats.IStats)
	attackerStats, _ := action.Actor.(stats.IStats)
	if targetStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction Update()")
		return 	
	}

	// check if already next to target
	if NextToEachOther(action.Actor, action.Target) {
		return
	} else {
		// move to target
		tx, ty := action.Target.GetPosition()
		action.Actor.SetPosition(tx, ty +1)
	}
}

func (action *AttackAction) CanBeExecuted() bool {
	targetStats, _ := action.Target.(stats.IStats)
	attackerStats, _ := action.Actor.(stats.IStats)
	if targetStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction CanBeExecuted()")
		return false // can not execute, invalid actor or target
	}

	targetHp := targetStats.GetStat(stats.HpCurrent)
	return targetHp > 0
}

func (action *AttackAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	targetStats, _ := action.Target.(stats.IStats)
	attackerStats, _ := action.Actor.(stats.IStats)
	if targetStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction Update()")
		return true	// action is complete we have invalid actor or target
	}

	// If target no longer next to us, our action is done
	if !NextToEachOther(action.Actor, action.Target) {
		return true
	}

	// we attack once per second
	action.Timer_s -= dt_s
	for action.Timer_s <= 0 {
		action.Timer_s += 1

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
			return true
		}
	}

	// harvesting is not complete so we return FALSE
	return false
}