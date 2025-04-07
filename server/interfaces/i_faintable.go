package interfaces

type IFaintable interface {
	Faint()
	DeltaStatAndFaintIfApplicable(string)
}