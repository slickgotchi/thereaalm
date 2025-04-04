package types

import (
	"fmt"
)

type Item struct {
    Name string
    Quantity int
}

type TradeOffer struct {
    SentItems []Item
    RequestedItems []Item
}

type IInventory interface {
    GetItems() []Item
    GetItemsMap() *map[string]int
    GetItemsExceptGold() []Item
	AddItem(name string, quantity int)
	RemoveItem(name string, quantity int) int 
    GetItemQuantity(name string) int
    ProposeTrade(responder IInventory, tradeOffer TradeOffer) bool
    RespondToTrade(initiator IInventory, tradeOffer TradeOffer) bool
	DisplayInventory()
}

// ITEM HOLDER
type Inventory struct {
	Items map[string]int
}

func NewInventory() *Inventory {
	return &Inventory{
		Items: make(map[string]int),
	}
}

func (inv *Inventory) GetItems() []Item {
    var items []Item

    for name, quantity := range inv.Items {
        items = append(items, Item{
            Name: name,
            Quantity: quantity,
        })
    }

    return items
}

func (inv *Inventory) GetItemsMap() *map[string]int {
    return &inv.Items
}

func (inv *Inventory) GetItemsExceptGold() []Item {
    var items []Item

    for name, quantity := range inv.Items {
        if name != "gold" {
            items = append(items, Item{
                Name: name,
                Quantity: quantity,
            })
        }
    }

    return items
}

// AddItem adds a quantity of an item to the inventory
func (inv *Inventory) AddItem(name string, quantity int) {
    inv.Items[name] += quantity
}

// RemoveItem removes a quantity of an item from the inventory
// and returns the actual amount removed (may be less than requested if inventory is insufficient)
// If the item quantity reaches zero, it will be removed from the map.
func (inv *Inventory) RemoveItem(name string, quantity int) int {
    if currentQuantity, exists := inv.Items[name]; exists {
        if currentQuantity >= quantity {
            inv.Items[name] -= quantity
            // If the quantity reaches zero, remove the item from the map
            if inv.Items[name] == 0 {
                delete(inv.Items, name)
            }
            return quantity
        }
        // If there are not enough items, remove all and return the remaining
        inv.Items[name] = 0
        delete(inv.Items, name) // Remove the item from the map once its quantity hits zero
        return currentQuantity
    }
    // If the item doesn't exist, return 0 (nothing was removed)
    return 0
}

func (inv *Inventory) GetItemQuantity(name string) int {
    if currentQuantity, exists := inv.Items[name]; exists {
        return currentQuantity
    }

    return 0
}

// TradeItem sends another ItemHolder a good in exchange for another
// returns TRUE if the trade was a SUCCESS
func (inv *Inventory) ProposeTrade(respondingItemHolder IInventory, 
    tradeOffer TradeOffer) bool {

    // check we have the items we want to send
    for _, item := range tradeOffer.SentItems {
        if inv.Items[item.Name] < item.Quantity {
            return false
        }
    }

    // check the offeree wants to / can do the offer
    isAccepted := respondingItemHolder.RespondToTrade(inv, tradeOffer)

    // update inventories
    if isAccepted {
        for _, item := range tradeOffer.SentItems {
            respondingItemHolder.AddItem(item.Name, item.Quantity)
            inv.RemoveItem(item.Name, item.Quantity)
        }
        for _, item := range tradeOffer.RequestedItems {
            inv.AddItem(item.Name, item.Quantity)
            respondingItemHolder.RemoveItem(item.Name, item.Quantity)
        }
    }
    
    return isAccepted
}

func (inv *Inventory) RespondToTrade(sendingItemHolder IInventory, 
    tradeOffer TradeOffer) bool {

    // check we have the requested items
    for _, item := range tradeOffer.RequestedItems {
        if inv.Items[item.Name] < item.Quantity {
            return false
        }
    }

    // accept anything by default if we have all requested items in stock
    return true
}

// For testing purposes: display the inventory
func (inv *Inventory) DisplayInventory() {
    fmt.Println("Inventory:")
    for name, quantity := range inv.Items {
        fmt.Printf("%s: %d\n", name, quantity)
    }
}