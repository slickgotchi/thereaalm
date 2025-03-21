package entity

import (
	"log"
	"thereaalm/types"

	"github.com/google/uuid"
)
type BerryBush struct {
	Entity
	Movable
	ItemHolder
	MaxBerries int
	// CurrentBerries int
	BerryTimer_s float64
}

func NewBerryBush(zoneId, x, y int) *BerryBush {
	newInventory := NewItemHolder()
	newInventory.Items["Berry"] = 50

	return &BerryBush{
        Entity: Entity{
            ID:   types.EntityUUID(uuid.New()),
            Type: "BerryBush",
        },
        Movable: Movable{
			ZoneID: zoneId,
            X: x,
            Y: y,
        },
		MaxBerries: 50,
		// CurrentBerries: 50,
		ItemHolder: *newInventory,
        // Actionable doesn't need to be explicitly initialized.
    }
}

func (b *BerryBush) GetType() string {
	return b.Type
}

func (b *BerryBush) Update(dt_s float64) {
	log.Printf("BerryBush at (%d, %d)", b.X, b.Y)
	b.DisplayInventory()

	respawnInterval := 20.0
	respawnAmount := 10

	// recharge the berries on the bush
	b.BerryTimer_s += dt_s
	if b.BerryTimer_s > respawnInterval {
		log.Println("Bush respawned 10 berry's")
		b.BerryTimer_s -= respawnInterval
		b.Items["Berry"] += respawnAmount
		if b.Items["Berry"] > b.MaxBerries {
			b.Items["Berry"] = b.MaxBerries
		}
	}
}

func (b *BerryBush) Harvest() (string, int) {
	harvestAmount := min(5, b.Items["Berry"])

	b.RemoveItem("Berry", harvestAmount)

	if harvestAmount > 0 {
		return "Berry", harvestAmount
	} else {
		return "", 0
	}

}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
