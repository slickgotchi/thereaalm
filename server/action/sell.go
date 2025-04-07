package action

import (
	"fmt"
	"log"
	"thereaalm/interfaces"
	"thereaalm/types"
	"thereaalm/utils"
	"time"
)

type SellAction struct {
	Action
	Duration_s float64
	Timer_s float64
	TradeType string
}

func NewSellAction(actor, target interfaces.IEntity, weighting float64,
	fallbackTargetSpec *types.TargetSpec) *SellAction {

	seller, _ := actor.(interfaces.IStats)
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

func (a *SellAction) IsValidTarget(potentialTarget interfaces.IEntity) bool {
	if potentialTarget == nil {
		return false
	}

	respondingItemHolder, _ := potentialTarget.(types.IInventory) 
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
	initiatingItemHolder, _ := potentialActor.(types.IInventory)

	// correct types?
	if initiatingItemHolder == nil {
		log.Printf("ERROR [%s]: Invalid actor, returning...", utils.GetFuncName())
		return false	// action is complete we have invalid actor or target
	}

	// has the initiator got any items?
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