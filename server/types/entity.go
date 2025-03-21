package types

import (
	"github.com/google/uuid"
)

type EntityUUID uuid.UUID

// IEntity is the core interface for all entities.
type IEntity interface {
    GetUUID() EntityUUID
    GetType() string
    Update(dt_s float64)
}

// IMovable is for entities that can move.
type IMovable interface {
    GetPosition() (int, int)
    SetPosition(x, y int)
}

// IActionSequence is for entities that can process actions.
type IActionSequence interface {
    QueueAction(a IAction)
    ProcessActions(dt_s float64)
}

type Item struct {
    Name string
    Quantity int
}

type IItemHolder interface {
    GetItems() []Item
	AddItem(name string, quantity int)
	RemoveItem(name string, quantity int) int 
    MakeTradeOffer(receiveItemHolder IItemHolder, sendItems []Item, receiveItems []Item) bool
    ReviewTradeOffer(receiveItemHolder IItemHolder, sendItems []Item, receiveItems []Item) bool
}

type IHarvestable interface {
    Harvest() (string, int)
}

