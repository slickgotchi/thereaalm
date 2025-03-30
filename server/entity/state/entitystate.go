package state

type EntityState string

const (
	Dead EntityState = "dead"
	Active      EntityState = "active"
	Idle EntityState = "idle"
	Roaming EntityState = "roaming"
	Repairing EntityState = "repairing"
	Building  EntityState = "building"
	Upgrading EntityState = "upgrading"
)