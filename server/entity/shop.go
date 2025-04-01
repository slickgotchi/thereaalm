package entity

import (
	// "log"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Shop struct {
    Entity
	types.Inventory
}

func NewShop(zoneId, x, y int) *Shop {
	// start show with gold
	itemHolder := types.NewInventory()
	itemHolder.Items["gold"] = 10000

    return &Shop{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "shop",
            X: x,
            Y: y,
        },
        Inventory: *itemHolder,
    }
}

func (s *Shop) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Inventory map[string]int `json:"inventory"`
	}{
		Name: s.Type,
		Description: "Buy and sell items from one convenient location",
		Inventory: s.Items,
	}
}

func (s *Shop) Update(dt_s float64) {

}
