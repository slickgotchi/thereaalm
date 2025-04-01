package resourceentity

import (
	"thereaalm/entity"

	"github.com/google/uuid"
)
type AlphaSlateBoulders struct {
	entity.Entity
}

func NewAlphaSlateBoulders(zoneId, x, y int) *AlphaSlateBoulders {
	return &AlphaSlateBoulders{
        Entity: entity.Entity{
            ID:   uuid.New(),
            Type: "alphaslateboulders",
			X: x,
			Y: y,
        },
    }
}

func (b *AlphaSlateBoulders) CanBeMined() bool {
	return true
}

func (b *AlphaSlateBoulders) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		Description string `json:"description"`
	}{
		Name: b.Type,
		Description: "A source of the rare AlphaSlate material",
	}
}

func (b *AlphaSlateBoulders) Update(dt_s float64) {
	// do nothing
}

func (b *AlphaSlateBoulders) Mine() (string, int) {
	return "alphaslate", 1
}
