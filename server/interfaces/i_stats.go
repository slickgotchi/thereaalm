package interfaces

type IStats interface {
    SetStat(name string, value int)
    GetStat(name string) int
    DeltaStat(name string, value int)
}