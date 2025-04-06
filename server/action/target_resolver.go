package action

import (
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
		// log.Println("checking for nearest entity of type ", a.FallbackTargetSpec.TargetType)
        if len(candidates) == 0 {
            return nil
        }

        var closest types.IEntity = nil
        minDist := 1000000
        ax, ay := a.Actor.GetPosition()

        for _, candidate := range candidates {
			// skip self
            if candidate.GetUUID() == a.Actor.GetUUID() {
                continue // Skip the actor itself
            }

			// skip invalid potential targets
			if !a.IsValidTarget(candidate) {
				continue
			}

			// find closest
            tx, ty := candidate.GetPosition()
            dist := AbsInt(ax - tx) + AbsInt(ay - ty)
            if dist < minDist {
                minDist = dist
                closest = candidate
            }
        }
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
			// skip self
			if candidate.GetUUID() == a.Actor.GetUUID() {
				continue // Skip the actor itself
			}

			// skip invalid potential targets
			if !a.IsValidTarget(candidate) {
				continue
			}

            if candidateStats, ok := candidate.(stats.IStats); ok {
                pulse := float64(candidateStats.GetStat(stats.Pulse))
                if pulse < minPulse {
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
            if amount > maxAmount {
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
                if inventory.GetItemQuantity(itemName) > 0 {
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
            if a.CanBeExecuted() {
                validCandidates = append(validCandidates, candidate)
            }
        }
        if len(validCandidates) == 0 {
            return nil
        }
        return validCandidates[rand.Intn(len(validCandidates))]
    },
}

// // isValidTarget checks if the target is compatible with the action type
// func (a *Action) isValidTarget(target types.IEntity, actionType string) bool {
//     switch actionType {
//     case "forage":
//         _, ok := target.(types.IForageable)
//         return ok
//     case "chop":
//         _, ok := target.(types.IChoppable)
//         return ok
//     case "mine":
//         _, ok := target.(types.IMineable)
//         return ok
//     case "attack":
//         targetStats, ok := target.(stats.IStats)
// 		ok = targetStats.GetStat(stats.Pulse) > 0
//         return ok && target.GetType() != a.Actor.GetType()
// 	case "rebuild":

//     case "maintain":
//         _, targetOk := target.(stats.IStats)
//         if !targetOk {
//             return false
//         }
//         _, actorOk := a.Actor.(types.IInventory)
//         if !actorOk {
//             return false
//         }
        
//     default:
//         return true
//     }

// 	return true
// }