package resourceentity

import (
	// "log"
	"thereaalm/entity"
	"thereaalm/types"
	"thereaalm/utils"
	"time"

	"github.com/google/uuid"
)
type FomoBerryBush struct {
	entity.Entity
	types.Inventory
	MaxBerries int
	RegrowInterval_s time.Duration
	RegrowAmount int
	TimeOfLastRegrow time.Duration
}

func NewFomoBerryBush(x, y int) *FomoBerryBush {
	newInventory := types.NewInventory()
	newInventory.Items["fomoberry"] = 50

	return &FomoBerryBush{
        Entity: entity.Entity{
            ID:   uuid.New(),
            Type: "fomoberrybush",
			X: x,
			Y: y,
        },
		MaxBerries: 50,
		RegrowInterval_s: 20 * time.Second,
		RegrowAmount: 10,
		Inventory: *newInventory,
    }
}

func (b *FomoBerryBush) CanBeForaged() bool {
	return b.Items["fomoberry"] > 0
}

func (b *FomoBerryBush) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Items interface{} `json:"items"`
	}{
		Name: b.Type,
		Description: "The berries from this bush restore Spark when eaten",
		Items: b.Items,
	}
}

func (b *FomoBerryBush) Update(dt_s float64) {
	if b.WorldManager.Since(b.TimeOfLastRegrow) >= b.RegrowInterval_s {
		b.Items["fomoberry"] += b.RegrowAmount
		if b.Items["fomoberry"] > b.MaxBerries {
			b.Items["fomoberry"] = b.MaxBerries
		}
	}
}

func (b *FomoBerryBush) Forage() (string, int) {
	harvestAmount := utils.Min(5, b.Items["fomoberry"])

	b.RemoveItem("fomoberry", harvestAmount)

	if harvestAmount > 0 {
		return "fomoberry", harvestAmount
	} else {
		return "", 0
	}

}