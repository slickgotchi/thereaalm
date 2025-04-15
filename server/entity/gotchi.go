package entity

import (
	// "log"
	"log"
	"strconv"
	"thereaalm/action"
	"thereaalm/entity/entitystate"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"
	"thereaalm/utils"
	"thereaalm/web3"
	"time"

	"github.com/google/uuid"
) 

type Gotchi struct {
    Entity
	action.ActionPlan
	types.Inventory
	Stats stattypes.Stats
	GotchiId string
	Name string
	SubgraphData web3.SubgraphGotchiData
	Personality []string
	Job string
	types.ActivityLog
	entitystate.State
	BuffMultiplier float64       // Buff-specific multiplier
    LastBuffCheck  time.Duration // Last buff check time
}

func NewGotchi(x, y int, subgraphGotchiData web3.SubgraphGotchiData) *Gotchi {
	// add item holder
	newItemHolder := types.NewInventory()

	// check for invalid subgraph data (requiring a default gotchi be used)
	if subgraphGotchiData.ID == "" {
		subgraphGotchiData = web3.DefaultSubgraphGotchiData
		log.Println("Using default gotchi data: ", subgraphGotchiData)
	}

	// get brs duration modifier
	// brsMultiplier := GetBRSMultiplier(subgraphGotchiData)
	// log.Print(brsMultiplier, int(400*brsMultiplier))

	// add some stats
	newStats := stattypes.NewStats()
	newStats.SetStat(stattypes.Ecto, 500)
	newStats.SetStat(stattypes.Spark, 500)
	newStats.SetStat(stattypes.Pulse, 500)
	newStats.SetStat(stattypes.MaxPulse, 1000)
	newStats.SetStat(stattypes.StakedGHST, 0)
	newStats.SetStat(stattypes.TreatTotal, 5000)

	// make new gotchi
	return &Gotchi{
        Entity: Entity{
            ID:   uuid.New(),
            Type: "gotchi",
			X: x,
			Y: y,
        },
        ActionPlan: action.ActionPlan{
			Actions: make([]interfaces.IAction, 0),
		},
		Inventory: *newItemHolder,
		Stats: *newStats,
		SubgraphData: subgraphGotchiData,
		Name: subgraphGotchiData.Name,
		GotchiId: subgraphGotchiData.ID,
		Personality: CreatePersonalityFromSubgraphData(subgraphGotchiData),
		State: entitystate.Active,
		BuffMultiplier: 1.0,
        LastBuffCheck:  0,
    }
}

func (g *Gotchi) GetSnapshotData() interface{} {
	return struct {
		Name string `json:"name"`
		UUID uuid.UUID `json:"uuid"`
		GotchiID  string `json:"gotchiId"`
		ZoneID int `json:"zoneId"`
		Description string `json:"description"`
		Stats interface{} `json:"stats"`
		Inventory interface{} `json:"inventory"`
		Personality interface{} `json:"personality"`
		Direction string `json:"direction"`
		ActivityLog interface{} `json:"activityLog"`
		ActionPlan interface{} `json:"actionPlan"`
		State entitystate.State `json:"state"`
		BuffMultiplier float64 `json:"buffmultiplier"`
		StakedGHST float64 `json:"stakedGhst"`
		TreatTotal float64 `json:"treatAmount"`
		Job string `json:"job"`
	}{
		Name: g.Name,
		UUID: g.ID,
		GotchiID:  g.GotchiId,
		ZoneID: g.GetZone().GetID(),
		Description: "The ethereal frens and sworn protectors of The Reaalm",
		Stats: g.Stats.StatMap,
		Inventory: g.Items,
		Personality: g.Personality,
		Direction: g.Direction,
		ActivityLog: g.ActivityLog.Entries,
		ActionPlan: g.ActionPlan.ToReporting(),
		State: g.State,
		BuffMultiplier: g.BuffMultiplier,
		Job: g.Job,
	}
}

func (g *Gotchi) Update(dt_s float64) {
	// calculate treat to add due to staking
	stakedGhst := g.Stats.GetStat(stattypes.StakedGHST)
	g.Stats.DeltaStat(stattypes.TreatTotal, stakedGhst / 86400 * dt_s)

	// if entity is dead don't do anything else
	if g.State == entitystate.Dead {
        return
    }

    // Check buffs periodically (e.g., every second)
    if g.WorldManager.Since(g.LastBuffCheck) >= time.Second {
        utils.UpdateBuffMultiplier(g, 16)
        g.LastBuffCheck = g.WorldManager.Now()
    }

	// process actions
    g.ProcessActions(dt_s)
}

// custom stat modification wrappers
func (e *Gotchi) SetStat(name string, value float64) {
	e.Stats.SetStat(name, value)
}

func (e *Gotchi) GetStat(name string) float64 {
	return e.Stats.GetStat(name)
}

func (e *Gotchi) DeltaStat(name string, value float64) {
	prev := e.Stats.GetStat(name)
	e.Stats.DeltaStat(name, value)
	newVal := e.Stats.GetStat(name)

	// ensure ecto, spark and pulse stay within range
	e.Stats.SetStat(stattypes.Ecto, utils.Clamp(e.Stats.GetStat(stattypes.Ecto), 0, 1000))
	e.Stats.SetStat(stattypes.Spark, utils.Clamp(e.Stats.GetStat(stattypes.Spark), 0, 1000))
	e.Stats.SetStat(stattypes.Pulse, utils.Clamp(e.Stats.GetStat(stattypes.Pulse), 0, 1000))

	// CUSTOM HOOK: handle ESP stats going below 0 (death)
	if (name == stattypes.Pulse || name == stattypes.Ecto || name == stattypes.Spark) && 
		newVal <= 0 && prev > 0 {

		// set gotchi state to dead
		e.State = entitystate.Dead

		// move gotchi to new location out of the way of entities
		currX, currY := e.GetPosition()
		newX, newY, found := 
			e.GetWorldManager().FindNearbyAvailablePosition(currX, currY, 7, 1)
		if found {
			// set direction to new position
			e.SetDirection("down")
	
			// Move the actor to the new position
			e.SetPosition(newX, newY)
		}
	}
}

// IBuffConsumer methods
func (g *Gotchi) GetEffectiveSpeedMultiplier() float64 {
    return g.BuffMultiplier
}

func (g *Gotchi) SetBuffMultiplier(multiplier float64) {
    g.BuffMultiplier = multiplier
}

func (g *Gotchi) GetLastBuffCheck() time.Duration {
    return g.LastBuffCheck
}

func (g *Gotchi) SetLastBuffCheck(t time.Duration) {
    g.LastBuffCheck = t
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