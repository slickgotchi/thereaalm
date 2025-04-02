package entity

import (
	// "log"
	"log"
	"strconv"
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
	types.ActivityLog
	ESP
}

type ESP struct {
	Ecto int
	Spark int
	Pulse int
}

func NewGotchi(zoneId, x, y int, subgraphGotchiData web3.SubgraphGotchiData) *Gotchi {
	// add item holder
	newItemHolder := types.NewInventory()

	// check for invalid subgraph data (requiring a default gotchi be used)
	if subgraphGotchiData.ID == "" {
		subgraphGotchiData = web3.DefaultSubgraphGotchiData
		log.Println("Using default gotchi data: ", subgraphGotchiData)
	}

	// get brs duration modifier
	brsMultiplier := GetBRSMultiplier(subgraphGotchiData)

	log.Print(brsMultiplier, int(400*brsMultiplier))

	// add some stats
	newStats := stats.NewStats()
	newStats.SetStat(stats.Ecto, 500)
	newStats.SetStat(stats.Spark, 500)
	newStats.SetStat(stats.Pulse, 500)

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
		ESP: ESP{Ecto: 500, Spark: 500, Pulse: 500},
    }
}

func (g *Gotchi) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		GotchiID  string `json:"gotchiId"`
		Description string `json:"description"`
		Stats interface{} `json:"stats"`
		Inventory interface{} `json:"inventory"`
		Personality interface{} `json:"personality"`
		Direction string `json:"direction"`
		ActivityLog interface{} `json:"activityLog"`
		ActionPlan interface{} `json:"actionPlan"`
	}{
		Name: g.Name,
		GotchiID:  g.GotchiId,
		Description: "The ethereal frens and sworn protectors of The Reaalm",
		Stats: g.Stats.StatMap,
		Inventory: g.Items,
		Personality: g.Personality,
		Direction: g.Direction,
		ActivityLog: g.ActivityLog.Entries,
		ActionPlan: g.ActionPlan.ToReporting(),
	}
}

func (g *Gotchi) Update(dt_s float64) {
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

func GetBRSMultiplier(subgraphData web3.SubgraphGotchiData) float64 {
	brs, _ := strconv.Atoi(subgraphData.WithSetsRarityScore)
	if brs < 500 {
		brs = 500
	}
	if brs > 1000 {
		brs = 1000
	}

	return float64(brs) / 500.0
}