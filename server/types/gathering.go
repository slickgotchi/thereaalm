package types

type IForageable interface {
    Forage() (string, int)
    CanBeForaged() bool
}

type IChoppable interface {
    Chop() (string, int)
    CanBeChopped() bool
}

type IHarvestable interface {
    Harvest() (string, int)
    CanBeHarvested() bool
}

type IMineable interface {
    Mine() (string, int)
    CanBeMined() bool
}