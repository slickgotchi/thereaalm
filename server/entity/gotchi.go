package entity

import (
	// "log"
	"log"
	"thereaalm/action"
	"thereaalm/stats"
	"thereaalm/types"
	"thereaalm/web3"

	"github.com/google/uuid"
)

type Gotchi struct {
    Entity
	action.ActionPlan
	types.Inventory
	stats.Stats
	GotchiId string
	Name string
	SubgraphData web3.SubgraphGotchiData
	Personality []string
}

// type Personality struct {
// 	Zen int `json:"zen"`
// 	Energetic int `json:"energetic"`
// 	Peaceful int `json:"peaceful"`
// 	Combative int `json:"combative"`
// 	Cuddly int `json:"cuddly"`
// 	Terrifying int `json:"terrifying"`
// 	Curious int `json:"curious"`
// 	Wise int `json:"wise"`
// 	Rugged int `json:"rugged"`
// 	Beautiful int `json:"beautiful"`
// 	Demonic int `json:"demonic"`
// 	Angelic int `json:"angelic"`
// }

func NewGotchi(zoneId, x, y int, subgraphGotchiData web3.SubgraphGotchiData) *Gotchi {
	// add item holder
	newItemHolder := types.NewInventory()

	// add some stats
	newStats := stats.NewStats()
	newStats.SetStat(stats.HpCurrent, 400)
	newStats.SetStat(stats.HpMax, 400)
	newStats.SetStat(stats.Attack, 5)
	newStats.SetStat(stats.HarvestDuration_s, 5)
	newStats.SetStat(stats.TradeDuration_s, 5)

	log.Println("NewGotchi: ", subgraphGotchiData)

	// make new gotchi
	return &Gotchi{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "gotchi",
			X: x,
			Y: y,
        },
        ActionPlan: action.ActionPlan{
			Actions: make([]types.IAction, 0),
		},
		Inventory: *newItemHolder,
		Stats: *newStats,
		SubgraphData: subgraphGotchiData,
		Name: subgraphGotchiData.Name,
		GotchiId: subgraphGotchiData.ID,
		Personality: CreatePersonalityFromSubgraphData(subgraphGotchiData),
    }
}

func (g *Gotchi) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		GotchiID  string `json:"gotchiId"`
		Stats interface{} `json:"stats"`
		Inventory interface{} `json:"inventory"`
		Personality interface{} `json:"personality"`
	}{
		Name: g.Name,
		GotchiID:  g.GotchiId,
		Stats: g.Stats.StatMap,
		Inventory: g.Items,
		Personality: g.Personality,
	}
}

func (g *Gotchi) Update(dt_s float64) {
	// log.Printf("Gotchi at (%d, %d)", g.X, g.Y)
	// g.DisplayInventory()

	// process our actions
	g.ProcessActions(dt_s);
}

func CreatePersonalityFromSubgraphData(subgraphData web3.SubgraphGotchiData) []string {
    traits := subgraphData.ModifiedNumericTraits
    if len(traits) < 4 {
        log.Printf("Warning: ModifiedNumericTraits has fewer than 4 elements: %v", traits)
        return []string{}
    }

    personality := make([]string, 0, 4) // Pre-allocate for 4 traits (NRG, AGG, SPK, BRN)

    // Helper function to determine the personality trait based on value
    addPersonalityTrait := func(value int, mythLow, rareLow, uncommonLow, commonLow, commonHigh, uncommonHigh, rareHigh, mythHigh string) {
        switch {
        case value <= 1: personality = append(personality, mythLow)
        case value >= 2 && value <= 9: personality = append(personality, rareLow)
        case value >= 10 && value <= 24: personality = append(personality, uncommonLow)
        case value >= 25 && value <= 49: personality = append(personality, commonLow)
        case value >= 50 && value <= 74: personality = append(personality, commonHigh)
        case value >= 75 && value <= 90: personality = append(personality, uncommonHigh)
        case value >= 91 && value <= 97: personality = append(personality, rareHigh)
        case value >= 98: personality = append(personality, mythHigh)
        default:
            log.Printf("Warning: Trait value %d is out of range (0-100)", value)
        }
    }

    // Energy (NRG) - Index 0
    addPersonalityTrait(traits[0],
        "Zen", "Serene", "Sleepy", "Calm",
        "Alert", "Energetic", "Hyper", "Turnt")

    // Aggressiveness (AGG) - Index 1
    addPersonalityTrait(traits[1],
        "Nonviolent", "Peaceful", "Forgiving", "Gentle",
        "Assertive", "Combative", "Warlike", "Based")

    // Spookiness (SPK) - Index 2
    addPersonalityTrait(traits[2],
        "Cuddly", "Impish", "Unnerving", "Shocking",
        "Scary", "Creepy", "Terrifying", "Ghastly")

    // Brain Size (BRN) - Index 3
    addPersonalityTrait(traits[3],
        "Glitchy", "Ditsy", "Quirky", "Basic",
        "Witty", "Genius", "Visionary", "Mindbender")

	// Eye Shape (EYS) - Index 4
	addPersonalityTrait(traits[4],
		"Ravishing", "Rugged", "Striking", "Cute",
		"Pretty", "Stunning", "Gorgeous", "Smokeshow")

	// Eye Color (EYC) - Index 5
	addPersonalityTrait(traits[5], 
		"Demonic", "Malicious", "Wicked", "Primal",
		"Earthy", "Bubbly", "Soulful", "Angelic")

    return personality
}