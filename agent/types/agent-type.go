package types

import "finalTaskLMS/globals"

type Agent struct {
	Tasks             []globals.PrimeEvaluation
	CountOfGoroutines int
	LimitOfGoroutines int
}
