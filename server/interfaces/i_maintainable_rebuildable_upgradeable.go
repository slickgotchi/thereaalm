package interfaces

type IMaintainable interface {
    Maintain(pulseRestored int)
    CanBeMaintained() bool
	GetMaxPulse() int
}

type IRebuildable interface {
    Rebuild(pulseRestored int)
    CanBeRebuilt() bool
	GetMaxPulse() int
}

type IUpgradeable interface {
    Upgrade() (string, int)
    CanBeUpgraded() bool
}