// thereaalm/entity/gotchi.go
package entities

import (
	"thereaalm/ai"
	"thereaalm/entity/component"
	"thereaalm/web3"

	"github.com/google/uuid"
)

type GotchiEntity struct {
    UUID       string
    Position *component.Position
    Velocity *component.Velocity
    Gotchi   *component.GotchiData
}

func NewGotchi(zoneID, x, y int, subgraphData web3.SubgraphGotchiData) *GotchiEntity {
    return &GotchiEntity{
        UUID: uuid.New().String(),
        Position: &component.Position{
            ZoneID: zoneID,
            X:      x,
            Y:      y,
        },
        Velocity: &component.Velocity{
            VX: 0,
            VY: 0,
        },
        Gotchi: &component.GotchiData{
            Mind:         ai.NewGotchiMind(),
            SubgraphData: subgraphData,
        },
    }
}

func (g *GotchiEntity) Update() {
    // Example update logic
    g.Position.X += g.Velocity.VX
    g.Position.Y += g.Velocity.VY
    
    g.Gotchi.Mind.Update() // Placeholder for AI logic
}

func (g *GotchiEntity) GetPosition() *component.Position {
    return g.Position
}