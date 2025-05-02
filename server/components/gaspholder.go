package components

type GASPHolder struct {
	Value int
}

func (gh *GASPHolder) AddGASP(amount int) {
	gh.Value += amount
}

func (gh *GASPHolder) GetGASP() int {
	return gh.Value
}

func (gh *GASPHolder) RemoveGASP(amount int) {
	gh.Value -= amount
}

func (gh *GASPHolder) TryRemoveGASP(amount int) bool {
	if gh.Value >= amount {
		gh.Value -= amount
		return true
	}

	return false
}