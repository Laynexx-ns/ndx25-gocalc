package handlers

import (
	"finalTaskLMS/internal/models"
	"finalTaskLMS/internal/services/agent/types"
	"testing"
)

func TestProcessTask(t *testing.T) {
	tests := []struct {
		name     string
		input    models.PrimeEvaluation
		expected float64
	}{
		{"Addition", models.PrimeEvaluation{Arg1: 1, Arg2: 2, Operation: "+"}, 3},
		{"Subtraction", models.PrimeEvaluation{Arg1: 10, Arg2: 2, Operation: "-"}, 8},
		{"Multiplication", models.PrimeEvaluation{Arg1: 3, Arg2: 5, Operation: "*"}, 15},
		{"Division", models.PrimeEvaluation{Arg1: 10, Arg2: 2, Operation: "/"}, 5},
	}

	agent := &types.Agent{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expr := test.input
			a := processTask(&expr, agent)

			if a != test.expected {
				t.Errorf("expected %f, got %f", test.expected, expr.Result)
			}
		})
	}
}
