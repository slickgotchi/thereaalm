package action

import (
	"log"
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

func (action *AttackAction) Start() {
		// check actor and target are of correct type
		targetStats, _ := action.Target.(types.IStats)
		attackerStats, _ := action.Actor.(types.IStats)
		if targetStats == nil || attackerStats == nil {
			log.Printf("Invalid IStats for actor or target in AttackAction Update()")
			return 	
		}
	
		// move to target
		tx, ty := action.Target.GetPosition()
		action.Actor.SetPosition(tx, ty +1)
}

func (action *AttackAction) CanBeExecuted() bool {
	targetStats, _ := action.Target.(types.IStats)
	attackerStats, _ := action.Actor.(types.IStats)
	if targetStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction CanBeExecuted()")
		return false // can not execute, invalid actor or target
	}

	targetHp, _ := targetStats.GetStatValue("hp")
	return targetHp > 0
}

func (action *AttackAction) Update(dt_s float64) bool {
	// check actor and target are of correct type
	targetStats, _ := action.Target.(types.IStats)
	attackerStats, _ := action.Actor.(types.IStats)
	if targetStats == nil || attackerStats == nil {
		log.Printf("Invalid IStats for actor or target in AttackAction Update()")
		return true	// action is complete we have invalid actor or target
	}

	// we attack once per second
	action.Timer_s -= dt_s
	for action.Timer_s <= 0 {
		action.Timer_s += 1

		attackDamage, ok := attackerStats.GetStatValue("attack"); 
		if !ok {
			log.Printf("attacker does not have attack stat")
			return true
		}

		targetStats.DeltaStatValue("hp", -attackDamage)
		log.Printf("entity did %d damage to anotherentity", attackDamage)
		newHp, _ := targetStats.GetStatValue("hp")
		log.Printf("newHp: %d", newHp);

		if newHp <= 0 {
			log.Println("Defeated enemy")
			return true
		}
	}

	// harvesting is not complete so we return FALSE
	return false
}