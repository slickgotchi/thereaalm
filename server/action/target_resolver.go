package action

import (
	"log"
	"math"
	"math/rand"
	"thereaalm/stats"
	"thereaalm/types"
)

// FallbackFunc defines the signature for a fallback function
type FallbackFunc func(*Action) types.IEntity

// fallbackHandlers maps criterion names to their corresponding fallback functions
var FallbackHandlers = map[string]FallbackFunc{
    "nearest": func(a *Action) types.IEntity {
        zone := a.Actor.GetZone()
        candidates := zone.GetEntitiesByType(a.FallbackTargetSpec.TargetType)
		log.Println("checking for nearest entity of type ", a.FallbackTargetSpec.TargetType)
        if len(candidates) == 0 {
            return nil
        }

        var closest types.IEntity
        minDist := math.Inf(1)
        ax, ay := a.Actor.GetPosition()

        for _, candidate := range candidates {
            if candidate.GetUUID() == a.Actor.GetUUID() {
                continue // Skip the actor itself
            }
            tx, ty := candidate.GetPosition()
            dist := math.Sqrt(math.Pow(float64(ax-tx), 2) + math.Pow(float64(ay-ty), 2))
            if dist < minDist && a.isValidTarget(candidate, a.Type) {
                minDist = dist
                closest = candidate
            }
        }
		log.Println("Found nearest ", closest)
        return closest
    },
    "lowest_pulse": func(a *Action) types.IEntity {
        zone := a.Actor.GetZone()
        candidates := zone.GetEntitiesByType(a.FallbackTargetSpec.TargetType)
        if len(candidates) == 0 {
            return nil
        }

        var lowest types.IEntity
        minPulse := math.Inf(1)

        for _, candidate := range candidates {
            if candidateStats, ok := candidate.(stats.IStats); ok {
                pulse := float64(candidateStats.GetStat(stats.Pulse))
                if pulse < minPulse && a.isValidTarget(candidate, a.Type) {
                    minPulse = pulse
                    lowest = candidate
                }
            }
        }
        return lowest
    },
    "resource_rich": func(a *Action) types.IEntity {
        zone := a.Actor.GetZone()
        candidates := zone.GetEntitiesByType(a.FallbackTargetSpec.TargetType)
        if len(candidates) == 0 {
            return nil
        }

        var richest types.IEntity
        maxAmount := 0

        for _, candidate := range candidates {
            amount := 0
            if forageable, ok := candidate.(types.IForageable); ok {
                _, amount = forageable.Forage()
            } else if mineable, ok := candidate.(types.IMineable); ok {
                _, amount = mineable.Mine()
            }
            if amount > maxAmount && a.isValidTarget(candidate, a.Type) {
                maxAmount = amount
                richest = candidate
            }
        }
        return richest
    },
    "has_specific_item": func(a *Action) types.IEntity {
        itemName, ok := a.FallbackTargetSpec.TargetValue.(string)
        if !ok {
            return nil
        }

        zone := a.Actor.GetZone()
        candidates := zone.GetEntitiesByType(a.FallbackTargetSpec.TargetType)
        if len(candidates) == 0 {
            return nil
        }

        for _, candidate := range candidates {
            if inventory, ok := candidate.(types.IInventory); ok {
                if inventory.GetItemQuantity(itemName) > 0 && a.isValidTarget(candidate, a.Type) {
                    return candidate
                }
            }
        }
        return nil
    },
    "random": func(a *Action) types.IEntity {
        zone := a.Actor.GetZone()
        candidates := zone.GetEntitiesByType(a.FallbackTargetSpec.TargetType)
        if len(candidates) == 0 {
            return nil
        }

        validCandidates := []types.IEntity{}
        for _, candidate := range candidates {
            if a.isValidTarget(candidate, a.Type) {
                validCandidates = append(validCandidates, candidate)
            }
        }
        if len(validCandidates) == 0 {
            return nil
        }
        return validCandidates[rand.Intn(len(validCandidates))]
    },
}

// isValidTarget checks if the target is compatible with the action type
func (a *Action) isValidTarget(target types.IEntity, actionType string) bool {
    switch actionType {
    case "forage":
        _, ok := target.(types.IForageable)
        return ok
    case "chop":
        _, ok := target.(types.IChoppable)
        return ok
    case "mine":
        _, ok := target.(types.IMineable)
        return ok
    case "attack":
        _, ok := target.(stats.IStats)
        return ok && target.GetType() != a.Actor.GetType()
    case "rebuild", "maintain":
        _, targetOk := target.(stats.IStats)
        if !targetOk {
            return false
        }
        _, actorOk := a.Actor.(types.IInventory)
        if !actorOk {
            return false
        }
        
    default:
        return true
    }

	return true
}