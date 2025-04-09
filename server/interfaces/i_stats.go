package interfaces

type IStats interface {
    SetStat(name string, value float64)
    GetStat(name string) float64
    DeltaStat(name string, value float64)
}