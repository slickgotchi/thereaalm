package tradeactions

import (
	"fmt"
	"log"
	"thereaalm/action"
	"thereaalm/interfaces"
	"thereaalm/types"
	"thereaalm/utils"
	"time"
)

type SellAction struct {
	action.Action
	Duration_s float64
	Timer_s float64
	// TradeType string
}

func NewSellAction(actor, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *SellAction {

	wm := actor.GetZone().GetWorldManager()
	
	a := &SellAction{
		Action: action.Action{
			Type: "sell",
			Weighting: weighting,
			Actor: actor,
			Target: target,
			WorldManager: wm,
		},
		Duration_s: 10,
		Timer_s: 10,
		// TradeType: "SellAllForGold",
	}

	a.SetFallbackTargetSpec(fallbackTargetSpec)

	return a
}

func (a *SellAction) IsValidTarget(potentialTarget interfaces.IEntity) bool {
	if potentialTarget == nil {
		return false
	}

	respondingItemHolder, _ := potentialTarget.(interfaces.IInventory) 
	if respondingItemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(potentialTarget) {
		return false
	}

	return true
}

func (a *SellAction) IsValidActor(potentialActor interfaces.IEntity) bool {
	// check actor and target are correct type
	initiatingItemHolder, _ := potentialActor.(interfaces.IInventory)

	// correct types?
	if initiatingItemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// has the initiator got any sellable items?
	if len(initiatingItemHolder.GetSellableItems()) <= 0 {
		return false
	}

	// ok can execute
	return true
}

func (a *SellAction) Start() {
	a.Timer_s = a.Duration_s

	a.TryMoveToTargetEntity(a.Target)
}

func (a *SellAction) Update(dt_s float64) bool {
	// check actor and target are correct type
	initiatingTrader, _ := a.Actor.(interfaces.ITrader)
	respondingTrader, _ := a.Target.(interfaces.ITrader) 
	if initiatingTrader == nil || respondingTrader == nil {
		log.Printf("Invalid item holders passed to SellAction Update()")
		return true
	}

	initiatingInventory, _ := a.Actor.(interfaces.IInventory)
	respondingInventory, _ := a.Target.(interfaces.IInventory)
	if initiatingInventory == nil || respondingInventory == nil {
		log.Printf("Invalid invetory holder passed to SellAction Upadte()")
		return true
	}
	
	a.Timer_s -= dt_s
	if a.Timer_s <= 0 {
		var isSuccess bool

		// create initial sell offer
		initialSellOffer, isSuccess := initiatingTrader.CreateSellOffer(respondingTrader)
		if !isSuccess {
			return true
		}

		// create counter sell offer
		counterSellOffer, isSuccess := respondingTrader.CounterSellOffer(initiatingTrader, initialSellOffer)
		if !isSuccess {
			return true
		}

		// TEMPORARY: for now we just accept the counterSellOffer

		// add GASP to seller, remove GASP from buyer
		initiatingTrader.AddGASP(counterSellOffer.GASP)
		respondingTrader.RemoveGASP(counterSellOffer.GASP)

		// remove items from seller inventory, add to buyer inventory
		itemCount := 0
		for _, itemToSell := range counterSellOffer.ItemsToSell {
			initiatingInventory.RemoveItem(itemToSell.Name, itemToSell.Quantity)
			respondingInventory.AddItem(itemToSell.Name, itemToSell.Quantity)
			itemCount += itemToSell.Quantity
		}

		// log activity
		if activityLog, ok := a.Actor.(types.IActivityLog); ok {
			var logEntry types.ActivityLogEntry
			if isSuccess {
				logEntry = types.ActivityLogEntry{
					Description: fmt.Sprintf("Trade accepted: Sold %d items for %d GASP", itemCount, counterSellOffer.GASP),
					LogTime:     time.Now(),
				}
			} else {
				logEntry = types.ActivityLogEntry{
					Description: "Trade rejected: No deal was made.",
					LogTime:     time.Now(),
				}
			}
			activityLog.NewLogEntry(logEntry)
		}

		return true
	}

	// did not complete so return false
	return false
}
