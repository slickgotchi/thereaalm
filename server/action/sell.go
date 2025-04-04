package action

import (
	"fmt"
	"log"
	"thereaalm/stats"
	"thereaalm/types"
	"time"
)

type SellAction struct {
	Action
	Duration_s float64
	Timer_s float64
	TradeType string
}

func NewSellAction(actor, target types.IEntity, weighting float64,
	fallbackTargetSpec *TargetSpec) *SellAction {

	seller, _ := actor.(stats.IStats)
	if seller == nil {
		log.Println("ERROR: Selling actor does not have IStats, returning...")
		return nil
	}

	sellDuration_s := 5	// we need to swap this for ESP stat calcs later
	if sellDuration_s <= 0 {
		log.Println("ERROR: Trading actor must have 'trade_duration_s' stat, returning...")
		return nil
	}
	
	a := &SellAction{
		Action: Action{
			Type: "sell",
			Weighting: weighting,
			Actor: actor,
			Target: target,
		},
		Duration_s: float64(sellDuration_s),
		Timer_s: float64(sellDuration_s),
		TradeType: "SellAllForGold",
	}

	a.SetFallbackTargetSpec(fallbackTargetSpec)

	return a
}

func (a *SellAction) CanBeExecuted() bool {
	// check current target validity and/or set a fallback if neccessary
	if !a.EnsureValidTarget() {
		return false
	}
	
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

func (a *SellAction) Start() {
	a.Timer_s = a.Duration_s

	a.TryMoveToTargetEntity(a.Target)
}



func (a *SellAction) Update(dt_s float64) bool {
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