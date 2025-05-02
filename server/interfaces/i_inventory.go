package interfaces

type Item struct {
    Name string
    Quantity int
}

type IInventory interface {
    GetSellableItems() []Item
    GetItemsMap() *map[string]int
    GetItemsExceptGold() []Item
	AddItem(name string, quantity int)
	RemoveItem(name string, quantity int) int 
    GetItemQuantity(name string) int
	DisplayInventory()
}