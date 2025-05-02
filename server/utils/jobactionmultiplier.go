package utils

import (
	"fmt"
	"thereaalm/interfaces"
)

func GetJobActionMultiplier(entity interfaces.IEntity, action string) (float64, error) {
	if entity.GetType() != "gotchi" {
		return 1, nil
	}

	gotchi, isValid := entity.(interfaces.IGotchi)
	if !isValid {
		return 1, nil
	}

	job := gotchi.GetJob()

    switch job {
    case "mercenary":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 7, nil
        case "roam": return 3, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "warden":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "thief":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "beastmaster":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "medic":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "merchant":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "crafter":
        switch action {
        case "maintain": return 2, nil
        case "rebuild": return 2, nil
        case "chop": return 2, nil
        case "forage": return 5, nil
        case "mine": return 2, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "farmer":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "minerjack":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 5, nil
        case "forage": return 1, nil
        case "mine": return 5, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "builder":
        switch action {
        case "maintain": return 5, nil
        case "rebuild": return 5, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "alchemist":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "explorer":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 3, nil
        case "roam": return 7, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    case "scout":
        switch action {
        case "maintain": return 1, nil
        case "rebuild": return 1, nil
        case "chop": return 1, nil
        case "forage": return 1, nil
        case "mine": return 1, nil
        case "attack": return 1, nil
        case "roam": return 1, nil
        case "sell": return 1, nil
        default: return 0, fmt.Errorf("action %s not found", action)
        }
    default:
        return 0, fmt.Errorf("job %s not found", job)
    }
}