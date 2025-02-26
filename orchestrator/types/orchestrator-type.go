package types

import (
	"finalTaskLMS/orchestrator/internal/models"
	"sync"
)

var once sync.Once

type Orchestrator struct {
	Mu    sync.Mutex
	Queue []models.Expressions
}

func NewOrchestrator() *Orchestrator {
	var o Orchestrator
	once.Do(func() {
		o = Orchestrator{
			Mu:    sync.Mutex{},
			Queue: []models.Expressions{},
		}
	})
	return &o
}
