package action

import (
	"fmt"
	"log"
	"thereaalm/stats"
	"thereaalm/types"
	"time"
)

type TradeAction struct {
	Action
	Duration_s float64
	Timer_s float64
	TradeType string
}

func NewTradeAction(actor, target types.IEntity, weighting float64, tradeType string) *TradeAction {
	trader, _ := actor.(stats.IStats)
	if trader == nil {
		log.Println("ERROR: Trading actor does not have IStats, returning...")
		return nil
	}

	traderDuration_s := 5	// we need to swap this for ESP stat calcs later
	if traderDuration_s <= 0 {
		log.Println("ERROR: Trading actor must have 'trade_duration_s' stat, returning...")
		return nil
	}
	
	return &TradeAction{
		Action: Action{
			Type: "trade",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
		Duration_s: float64(traderDuration_s),
		Timer_s: float64(traderDuration_s),
		TradeType: tradeType,
	}
}

func (a *TradeAction) CanBeExecuted() bool {
	// check actor and target are correct type
	respondingItemHolder, _ := a.Target.(types.IInventory) 
	initiatingItemHolder, _ := a.Actor.(types.IInventory)

	// correct types?
	if respondingItemHolder == nil || initiatingItemHolder == nil {
		log.Printf("Invalid item holders passed to SellAction Update()")
		return false
	}

	// can move to target?
	if !a.CanMoveToTargetEntity(a.Target) {
		return false
	}

	// has items?
	if len(initiatingItemHolder.GetItemsExceptGold()) <= 0 {
		return false
	}

	// ok can execute
	return true
}

func (a *TradeAction) Start() {
	a.Timer_s = a.Duration_s

	a.TryMoveToTargetEntity(a.Target)
}



func (a *TradeAction) Update(dt_s float64) bool {
	// check actor and target are correct type
	respondingItemHolder, _ := a.Target.(types.IInventory) 
	initiatingItemHolder, _ := a.Actor.(types.IInventory)
	if respondingItemHolder == nil || initiatingItemHolder == nil {
		log.Printf("Invalid item holders passed to SellAction Update()")
		return true
	}
	
	a.Timer_s -= dt_s
	if a.Timer_s <= 0 {
		// this is where we iterate over different trade types OR
		// we insert custom logic from the holders that dictate
		// what they have for sale, what price they want to sell/buy at etc.
		if a.TradeType == "SellAllForGold" {
			// add up all items that aren't gold
			count := 0
			allInitiatorItems := initiatingItemHolder.GetItems()
			var filteredInitiatorItems []types.Item
			for _, item := range allInitiatorItems {
				if item.Name != "gold" {
					count += item.Quantity
					filteredInitiatorItems = append(filteredInitiatorItems, item)
				} 
			}

			var requestedItems []types.Item
			requestedItems = append(requestedItems, types.Item{
				Name: "gold",
				Quantity: count * 5,
			})

			tradeOffer := types.TradeOffer{
				SentItems: filteredInitiatorItems,
				RequestedItems: requestedItems,
			}

			// make the trade offer
			isAccepted := initiatingItemHolder.ProposeTrade(respondingItemHolder, tradeOffer)

			// Check if the actor has an activity log
			if activityLog, ok := a.Actor.(types.IActivityLog); ok {
				var logEntry types.ActivityLogEntry
				if isAccepted {
					logEntry = types.ActivityLogEntry{
						Description: fmt.Sprintf("Trade accepted: Sold %d items for %d gold", count, count*5),
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
	}

	// did not complete so return false
	return false
}