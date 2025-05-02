package interfaces

type IMaintainable interface {
    Maintain(pulseRestored float64)
    CanBeMaintained() bool
}

type IRebuildable interface {
    Rebuild(pulseRestored float64)
    CanBeRebuilt() bool
}

type IUpgradeable interface {
    Upgrade() (string, int)
    CanBeUpgraded() bool
}