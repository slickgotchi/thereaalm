package ai

import (
	"thereaalm/taskpool"
)

type DecisionEngine struct {
	Pool *taskpool.Pool
}

func NewDecisionEngine(pool *taskpool.Pool) *DecisionEngine {
	return &DecisionEngine{Pool: pool}
}
