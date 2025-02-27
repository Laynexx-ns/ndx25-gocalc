package handlers

import (
	"bytes"
	"encoding/json"
	"finalTaskLMS/agent/types"
	"finalTaskLMS/globals"
	"log"
	"net/http"
	"time"
)

func CycleTask(a *types.Agent) {
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			task, err := getTask()
			a.Tasks = append(a.Tasks, task)
			if err != nil {
				log.Println("Error getting task:", err)
				continue
			}

		}
	}
}

func getTask() (globals.PrimeEvaluation, error) {
	client := http.Client{}
	resp, err := client.Get("http://localhost:8080/internal/task")
	if err != nil {
		return globals.PrimeEvaluation{}, err
	}
	defer resp.Body.Close()

	var expression globals.PrimeEvaluation
	err = json.NewDecoder(resp.Body).Decode(&expression)
	if err != nil {
		return globals.PrimeEvaluation{}, err
	}

	return expression, nil
}

func processTask(expression *globals.PrimeEvaluation, a *types.Agent) {

	result := 0.0
	switch expression.Operation {
	case '+':
		result = expression.Arg1 + expression.Arg2
	case '-':
		result = expression.Arg1 - expression.Arg2
	case '*':
		result = expression.Arg1 * expression.Arg2
	case '/':
		if expression.Arg2 != 0 {
			result = expression.Arg1 / expression.Arg2
		} else {
			log.Println("Division by zero")
			return
		}
	default:
		log.Println("Unknown operation")
		return
	}

	log.Printf("Task %d calculated: %f", expression.Id, result)

	response := globals.PrimeEvaluationResponse{
		ParentID:      expression.ParentId,
		Id:            expression.Id,
		Arg1:          int(expression.Arg1),
		Arg2:          int(expression.Arg2),
		Operation:     expression.Operation,
		Result:        result,
		OperationTime: 2,
	}

	p, err := json.Marshal(response)
	if err != nil {
		log.Println("invalid expression:", err)
		return
	}

	client := http.Client{}
	_, err = client.Post("http://localhost:8080/internal/task",
		"application/json", bytes.NewReader(p))
	if err != nil {
		log.Println("can't send request to orchestrator")
	} else {
		log.Printf("evaluated expression was successfully send  - %v", p)
	}
}
