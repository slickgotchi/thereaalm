package interfaces

type IGASPHolder interface {
	AddGASP(amount int)
	GetGASP() int
	RemoveGASP() bool
	TryRemoveGASP(amount int) bool
}