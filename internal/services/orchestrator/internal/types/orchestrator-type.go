package types

import (
	"ndx/internal/models"
	"sync"
)

var once sync.Once

type Orchestrator struct {
	Mu              sync.Mutex
	Queue           []models.PrimeEvaluation
	SentEvaluations []models.PrimeEvaluation
	Expressions     []models.Expressions
	Subs            map[int]chan struct{}
	Chans           map[int]chan float64
	Errchans        map[int]chan error
}

func NewOrchestrator() *Orchestrator {
	var o Orchestrator
	once.Do(func() {
		o = Orchestrator{
			Mu:          sync.Mutex{},
			Queue:       []models.PrimeEvaluation{},
			Expressions: []models.Expressions{},
			Subs:        make(map[int]chan struct{}),
			Chans:       make(map[int]chan float64, 1),
			Errchans:    make(map[int]chan error, 1),
		}
	})
	return &o
}
