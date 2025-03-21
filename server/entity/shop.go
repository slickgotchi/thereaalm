package entity

import (
	"log"
	"thereaalm/types"

	"github.com/google/uuid"
)

type Shop struct {
    Entity
    Movable
	ItemHolder
}

func NewShop(zoneId, x, y int) *Shop {
	// start show with gold
	itemHolder := NewItemHolder()
	itemHolder.Items["Gold"] = 10000

    return &Shop{
        Entity: Entity{
            ID:   types.EntityUUID(uuid.New()),
            Type: "Shop",
        },
        Movable: Movable{
			ZoneID: zoneId,
            X: x,
            Y: y,
        },
        ItemHolder: *itemHolder,
    }
}

func (s *Shop) Update(dt_s float64) {
	log.Printf("Shop at (%d, %d)", s.X, s.Y)
	s.DisplayInventory()
}
