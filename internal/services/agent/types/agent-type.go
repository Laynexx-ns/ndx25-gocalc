package types

import (
	"finalTaskLMS/internal/models"
)

type Agent struct {
	Tasks             []models.PrimeEvaluation
	CountOfGoroutines int
	LimitOfGoroutines int
}
