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

func NewTradeAction(actor, target types.IEntity, duration_s float64, tradeType string) *TradeAction {
	return &TradeAction{
		Action: Action{
			Type: "Trade",
			IsActive: false,
			Actor: actor,
			Target: target,
		},
		Duration_s: duration_s,
		Timer_s: 0,
		TradeType: tradeType,
	}
}

func (sa *TradeAction) Update(dt_s float64) bool {
	// check actor and target are correct type
	buyingItemHolder, _ := sa.Target.(types.IItemHolder) 
	sellingItemHolder, _ := sa.Actor.(types.IItemHolder)
	if buyingItemHolder == nil || sellingItemHolder == nil {
		log.Printf("Invalid item holders passed to SellAction Update()")
		return true
	}

	sa.Timer_s -= dt_s
	if sa.Timer_s <= 0 {
		// this is where we iterate over different trade types OR
		// we insert custom logic from the holders that dictate
		// what they have for sale, what price they want to sell/buy at etc.
		if sa.TradeType == "SellAllForGold" {
			// add up all items that aren't gold
			count := 0
			allSellerItems := sellingItemHolder.GetItems()
			var filteredSellerItems []types.Item
			for _, item := range allSellerItems {
				if item.Name != "Gold" {
					count += item.Quantity
					filteredSellerItems = append(filteredSellerItems, item)
				} 
			}

			var receiveItems []types.Item
			receiveItems = append(receiveItems, types.Item{
				Name: "Gold",
				Quantity: count * 5,
			})

			// make the trade offer
			sellingItemHolder.MakeTradeOffer(buyingItemHolder, filteredSellerItems, receiveItems)
			return true
		}
	}

	// did not complete so return false
	return false
}