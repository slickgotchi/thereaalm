package entity

import (
	"fmt"
	"thereaalm/types"
)

type Entity struct {
	ID types.EntityUUID
	Type string
}

func (e *Entity) GetUUID() types.EntityUUID { return e.ID }
func (e *Entity) GetType() string          { return e.Type }
func (e *Entity) Update(dt_s float64) {}


type Movable struct {
	ZoneID int
	X int
	Y int
}

func (m *Movable) GetZoneID() int {
	return m.ZoneID
}

func (m *Movable) GetPosition() (int, int) {
	return m.X, m.Y
}

func (m *Movable) SetPosition(x, y int) {
	m.X = x
	m.Y = y
}

type ActionSequence struct {
    Actions []types.IAction
}

func (a *ActionSequence) QueueAction(action types.IAction) {
    a.Actions = append(a.Actions, action)
}

func (a *ActionSequence) ProcessActions(dt_s float64) {
    for _, action := range a.Actions {
        action.Execute(action.GetActor(), action.GetTarget()) // You can replace `nil` with actual actor/target when needed
    }
    a.Actions = nil // Clear actions after processing, if needed
}

type ItemHolder struct {
	Items map[string]int
}

func NewItemHolder() *ItemHolder {
	return &ItemHolder{
		Items: make(map[string]int),
	}
}

func (ih *ItemHolder) GetItems() []types.Item {
    var items []types.Item

    for name, quantity := range ih.Items {
        items = append(items, types.Item{
            Name: name,
            Quantity: quantity,
        })
    }

    return items
}

// AddItem adds a quantity of an item to the inventory
func (ih *ItemHolder) AddItem(name string, quantity int) {
    ih.Items[name] += quantity
}

// RemoveItem removes a quantity of an item from the inventory
// and returns the actual amount removed (may be less than requested if inventory is insufficient)
// If the item quantity reaches zero, it will be removed from the map.
func (ih *ItemHolder) RemoveItem(name string, quantity int) int {
    if currentQuantity, exists := ih.Items[name]; exists {
        if currentQuantity >= quantity {
            ih.Items[name] -= quantity
            // If the quantity reaches zero, remove the item from the map
            if ih.Items[name] == 0 {
                delete(ih.Items, name)
            }
            return quantity
        }
        // If there are not enough items, remove all and return the remaining
        ih.Items[name] = 0
        delete(ih.Items, name) // Remove the item from the map once its quantity hits zero
        return currentQuantity
    }
    // If the item doesn't exist, return 0 (nothing was removed)
    return 0
}

// TradeItem sends another ItemHolder a good in exchange for another
// returns TRUE if the trade was a SUCCESS
func (ih *ItemHolder) MakeTradeOffer(receivingItemHolder types.IItemHolder, 
    sendItems []types.Item, receiveItems []types.Item) bool {

    // check we have the items we want to send
    for _, item := range sendItems {
        if ih.Items[item.Name] < item.Quantity {
            return false
        }
    }

    // check the offeree wants to / can do the offer
    isAccepted := receivingItemHolder.ReviewTradeOffer(ih, sendItems, receiveItems)

    // update inventories
    if isAccepted {
        for _, item := range sendItems {
            receivingItemHolder.AddItem(item.Name, item.Quantity)
            ih.RemoveItem(item.Name, item.Quantity)
        }
        for _, item := range receiveItems {
            ih.AddItem(item.Name, item.Quantity)
            receivingItemHolder.RemoveItem(item.Name, item.Quantity)
        }
    }
    
    return isAccepted
}

func (ih *ItemHolder) ReviewTradeOffer(sendingItemHolder types.IItemHolder, 
    sendItems []types.Item, receiveItems []types.Item) bool {

    // check we have the requested items
    for _, item := range receiveItems {
        if ih.Items[item.Name] < item.Quantity {
            return false
        }
    }

    // accept anything by default if we have all requested items in stock
    return true
}

// For testing purposes: display the inventory
func (inv *ItemHolder) DisplayInventory() {
    fmt.Println("Inventory:")
    for name, quantity := range inv.Items {
        fmt.Printf("%s: %d\n", name, quantity)
    }
}