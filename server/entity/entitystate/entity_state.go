package entitystate

type State string

const (
	Dead State = "dead"
	Active      State = "active"
	Idle State = "idle"
	Roaming State = "roaming"
	Repairing State = "repairing"
	Building  State = "building"
	Upgrading State = "upgrading"
	Regrowing State = "regrowing"
)

func (es State) GetState() State {
	return es
}

type IEntityState interface {
	GetState() State
}
