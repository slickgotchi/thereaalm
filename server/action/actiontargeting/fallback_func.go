package actiontargeting

import (
	"math"
	"math/rand"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/types"
	"thereaalm/utils"
)

type FallbackCriteria string

const (
	FallbackNearest      FallbackCriteria = "nearest"
	FallbackLowestPulse  FallbackCriteria = "lowest_pulse"
	FallbackResourceRich FallbackCriteria = "resource_rich"
	FallbackHasItem      FallbackCriteria = "has_specific_item"
	FallbackRandom       FallbackCriteria = "random"
)

var FallbackHandlers = map[FallbackCriteria]func(a interfaces.IAction) interfaces.IEntity{
	FallbackNearest:      fallbackNearest,
	FallbackLowestPulse:  fallbackLowestPulse,
	FallbackResourceRich: fallbackResourceRich,
	FallbackHasItem:      fallbackHasSpecificItem,
	FallbackRandom:       fallbackRandom,
}

// --- Shared Helpers ---

func filterValidCandidates(a interfaces.IAction) []interfaces.IEntity {
	candidates := a.GetActor().GetZone().GetEntitiesByType(a.GetFallbackTargetSpec().TargetType)
	valid := make([]interfaces.IEntity, 0, len(candidates))
	for _, c := range candidates {
		if c.GetUUID() != a.GetActor().GetUUID() && a.IsValidTarget(c) {
			valid = append(valid, c)
		}
	}
	return valid
}

// --- Fallback Implementations ---

func fallbackNearest(a interfaces.IAction) interfaces.IEntity {
	valid := filterValidCandidates(a)
	if len(valid) == 0 {
		return nil
	}

	ax, ay := a.GetActor().GetPosition()
	var closest interfaces.IEntity
	minDist := 32

	for _, candidate := range valid {
		tx, ty := candidate.GetPosition()
		dist := utils.Abs(ax - tx) + utils.Abs(ay - ty)
		if dist < minDist {
			minDist = dist
			closest = candidate
		}
	}
	return closest
}

func fallbackLowestPulse(a interfaces.IAction) interfaces.IEntity {
	valid := filterValidCandidates(a)
	var lowest interfaces.IEntity
	minPulse := math.Inf(1)

	for _, candidate := range valid {
		if statsEntity, ok := candidate.(interfaces.IStats); ok {
			pulse := float64(statsEntity.GetStat(stattypes.Pulse))
			if pulse < minPulse {
				minPulse = pulse
				lowest = candidate
			}
		}
	}
	return lowest
}

func fallbackResourceRich(a interfaces.IAction) interfaces.IEntity {
	valid := filterValidCandidates(a)
	var richest interfaces.IEntity
	maxAmount := 0

	for _, candidate := range valid {
		amount := 0
		if f, ok := candidate.(types.IForageable); ok {
			_, amount = f.Forage()
		} else if m, ok := candidate.(types.IMineable); ok {
			_, amount = m.Mine()
		}
		if amount > maxAmount {
			maxAmount = amount
			richest = candidate
		}
	}
	return richest
}

func fallbackHasSpecificItem(a interfaces.IAction) interfaces.IEntity {
	itemName, ok := a.GetFallbackTargetSpec().TargetValue.(string)
	if !ok {
		return nil
	}

	candidates := a.GetActor().GetZone().GetEntitiesByType(a.GetFallbackTargetSpec().TargetType)
	for _, candidate := range candidates {
		if inventory, ok := candidate.(interfaces.IInventory); ok {
			if inventory.GetItemQuantity(itemName) > 0 {
				return candidate
			}
		}
	}
	return nil
}

func fallbackRandom(a interfaces.IAction) interfaces.IEntity {
	valid := filterValidCandidates(a)
	if len(valid) == 0 {
		return nil
	}
	return valid[rand.Intn(len(valid))]
}

func ResolveFallbackTarget(a interfaces.IAction) interfaces.IEntity {
	// ensure we have fallback spec
	if a.GetFallbackTargetSpec() != nil {
		criterion := FallbackCriteria(a.GetFallbackTargetSpec().TargetCriterion)
		if handler, ok := FallbackHandlers[criterion]; ok {
			return handler(a)
		}
	}

	return nil
}
