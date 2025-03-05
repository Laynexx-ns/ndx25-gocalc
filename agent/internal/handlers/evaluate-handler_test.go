package handlers

import (
	"finalTaskLMS/agent/types"
	"finalTaskLMS/globals"
	"testing"
)

func TestProcessTask(t *testing.T) {
	tests := []struct {
		name     string
		input    globals.PrimeEvaluation
		expected float64
	}{
		{"Addition", globals.PrimeEvaluation{Arg1: 1, Arg2: 2, Operation: "+"}, 3},
		{"Subtraction", globals.PrimeEvaluation{Arg1: 10, Arg2: 2, Operation: "-"}, 8},
		{"Multiplication", globals.PrimeEvaluation{Arg1: 3, Arg2: 5, Operation: "*"}, 15},
		{"Division", globals.PrimeEvaluation{Arg1: 10, Arg2: 2, Operation: "/"}, 5},
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
