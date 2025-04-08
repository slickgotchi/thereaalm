package interfaces

type IMaintainable interface {
    Maintain(pulseRestored int)
    CanBeMaintained() bool
}

type IRebuildable interface {
    Rebuild(pulseRestored int)
    CanBeRebuilt() bool
}

type IUpgradeable interface {
    Upgrade() (string, int)
    CanBeUpgraded() bool
}