package components

import (
	"fmt"
	"thereaalm/interfaces"
)

// ITEM HOLDER
type Inventory struct {
	Items map[string]int
}

func NewInventory() *Inventory {
	return &Inventory{
		Items: make(map[string]int),
	}
}

func (inv *Inventory) GetSellableItems() []interfaces.Item {
	// base implementation returns all items

    var items []interfaces.Item

    for name, quantity := range inv.Items {
        items = append(items, interfaces.Item{
            Name: name,
            Quantity: quantity,
        })
    }

    return items
}

func (inv *Inventory) GetItemsMap() *map[string]int {
    return &inv.Items
}

func (inv *Inventory) GetItemsExceptGold() []interfaces.Item {
    var items []interfaces.Item

    for name, quantity := range inv.Items {
        if name != "gold" {
            items = append(items, interfaces.Item{
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



// For testing purposes: display the inventory
func (inv *Inventory) DisplayInventory() {
    fmt.Println("Inventory:")
    for name, quantity := range inv.Items {
        fmt.Printf("%s: %d\n", name, quantity)
    }
}