package action

type TargetSpec struct {
    TargetType string
    TargetCriterion string      // e.g., "nearest", "resource_rich"
    TargetValue     interface{} // Optional parameter (e.g., "kekwood" for "has_item")
    SelfCriterion   string      // e.g., "min_pulse", "has_item"
    SelfValue       interface{} // e.g., 300 for "min_pulse"
}