package action

import (
	"log"
	"thereaalm/types"
)

type TradeAction struct {
	Action
	Duration_s float64
	Timer_s float64
	TradeType string
}

func NewTradeAction(actor, target types.IEntity, tradeType string) *TradeAction {
	trader, _ := actor.(types.IStats)
	if trader == nil {
		log.Println("ERROR: Trading actor does not have IStats, returning...")
		return nil
	}

	traderDuration_s, ok := trader.GetStatValue("trade_duration_s")
	if !ok {
		log.Println("ERROR: Trading actor must have 'trade_duration_s' stat, returning...")
		return nil
	}
	
	
	return &TradeAction{
		Action: Action{
			Type: "trade",
			IsStarted: false,
			Actor: actor,
			Target: target,
		},
		Duration_s: float64(traderDuration_s),
		Timer_s: 0,
		TradeType: tradeType,
	}
}

func (action *TradeAction) Update(dt_s float64) bool {
	// check actor and target are correct type
	respondingItemHolder, _ := action.Target.(types.IInventory) 
	initiatingItemHolder, _ := action.Actor.(types.IInventory)
	if respondingItemHolder == nil || initiatingItemHolder == nil {
		log.Printf("Invalid item holders passed to SellAction Update()")
		return true
	}

	// if first time, move to target
	if (!action.IsStarted) {
		action.IsStarted = true

		tx, ty := action.Target.GetPosition()
		action.Actor.SetPosition(tx, ty +1)
	}

	action.Timer_s -= dt_s
	if action.Timer_s <= 0 {
		// this is where we iterate over different trade types OR
		// we insert custom logic from the holders that dictate
		// what they have for sale, what price they want to sell/buy at etc.
		if action.TradeType == "SellAllForGold" {
			// add up all items that aren't gold
			count := 0
			allInitiatorItems := initiatingItemHolder.GetItems()
			var filteredInitiatorItems []types.Item
			for _, item := range allInitiatorItems {
				if item.Name != "Gold" {
					count += item.Quantity
					filteredInitiatorItems = append(filteredInitiatorItems, item)
				} 
			}

			var requestedItems []types.Item
			requestedItems = append(requestedItems, types.Item{
				Name: "Gold",
				Quantity: count * 5,
			})

			tradeOffer := types.TradeOffer{
				SentItems: filteredInitiatorItems,
				RequestedItems: requestedItems,
			}

			// make the trade offer
			initiatingItemHolder.ProposeTrade(respondingItemHolder, tradeOffer)
			return true
		}
	}

	// did not complete so return false
	return false
}