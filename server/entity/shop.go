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
		Inventory map[string]int `json:"inventory"`
	}{
		Inventory: s.Items,
	}
}

func (s *Shop) Update(dt_s float64) {
	// log.Printf("Shop at (%d, %d)", s.X, s.Y)
	// s.DisplayInventory()
}
