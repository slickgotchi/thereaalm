package resourceentity

import (
	// "log"
	"thereaalm/components"
	"thereaalm/entity"
	"thereaalm/entity/entitystate"
	"thereaalm/utils"
	"time"

	"github.com/google/uuid"
)
type KekWoodTree struct {
	entity.Entity
	components.Inventory
	MaxWood int
	RegrowDuration_s time.Duration
	TimeOfDepletion time.Duration
	State entitystate.State
}

func NewKekWoodTree(x, y int) *KekWoodTree {
	newInventory := components.NewInventory()
	newInventory.Items["kekwood"] = 100

	return &KekWoodTree{
        Entity: entity.Entity{
            ID:   uuid.New(),
            Type: "kekwoodtree",
			X: x,
			Y: y,
        },
		MaxWood: 100,
		Inventory: *newInventory,
		State: entitystate.Active,
    }
}

func (b *KekWoodTree) CanBeChopped() bool {
	return b.State == entitystate.Active
}

func (b *KekWoodTree) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Items interface{} `json:"items"`
	}{
		Name: b.Type,
		Description: "The berries from this bush fetch a pretty penny with shops",
		Items: b.Items,
	}
}

func (b *KekWoodTree) Update(dt_s float64) {

	// check if we're at 0 
	if b.State == entitystate.Active {
		if b.Items["kekwood"] <= 0 {
			b.State = entitystate.Regrowing
			b.TimeOfDepletion = b.WorldManager.Now()
		}
	}

	if b.State == entitystate.Regrowing {
		if b.WorldManager.Since(b.TimeOfDepletion) >= b.RegrowDuration_s {
			b.Items["kekwood"] = b.MaxWood
			b.State = entitystate.Active
		}
	}
}

func (b *KekWoodTree) Chop() (string, int) {
	chopAmount := utils.Min(5, b.Items["kekwood"])

	b.RemoveItem("kekwood", chopAmount)

	if chopAmount > 0 {
		return "kekwood", chopAmount
	} else {
		return "", 0
	}

}
