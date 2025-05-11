package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"ndx/internal/models"
	"ndx/internal/services/agent/internal"
	"ndx/pkg/config"

	"log"
	"net/http"
	"time"
)

type EvaluateHandler struct {
	db     *sql.DB
	config config.Config
	Agent  *internal.Agent
}

func NewEvaluateHandler(db *sql.DB, c config.Config) *EvaluateHandler {
	return &EvaluateHandler{
		db:     db,
		config: c,
	}
}

func (eh *EvaluateHandler) CycleTask() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			task, err := eh.getTask()
			eh.Agent.Tasks = append(eh.Agent.Tasks, task)
			if err != nil {
				log.Println("Error getting task:", err)
				continue
			}
			eh.processTask(&task)

		}
	}
}

func (eh *EvaluateHandler) getTask() (models.PrimeEvaluation, error) {
	client := http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://%s:%d/internal/task", eh.config.OrchestratorConf.Host, eh.config.OrchestratorConf.Port))
	if err != nil {
		return models.PrimeEvaluation{}, err
	}
	defer resp.Body.Close()

	var expression models.PrimeEvaluation
	err = json.NewDecoder(resp.Body).Decode(&expression)
	if err != nil {
		return models.PrimeEvaluation{}, err
	}

	return expression, nil
}

func (eh *EvaluateHandler) processTask(expression *models.PrimeEvaluation) float64 {

	result := 0.0
	operationTime := 0
	hasErr := false
	switch expression.Operation {
	case "^":
		result = math.Pow(expression.Arg1, expression.Arg2)
		operationTime = eh.config.AgentConf.TimeConf.MultiplicationTime
	case "+":
		result = expression.Arg1 + expression.Arg2
		operationTime = eh.config.AgentConf.TimeConf.AdditionTime
	case "-":
		result = expression.Arg1 - expression.Arg2
		operationTime = eh.config.AgentConf.TimeConf.SubtractionTime
	case "*":
		result = expression.Arg1 * expression.Arg2
		operationTime = eh.config.AgentConf.TimeConf.MultiplicationTime
	case "/":
		if expression.Arg2 != 0 {
			operationTime = eh.config.AgentConf.TimeConf.DivisionTime
			result = expression.Arg1 / expression.Arg2
		} else {
			log.Println("Division by zero")
			hasErr = true
			return 0
		}
	default:
		log.Println("Unknown operation", expression)
		hasErr = true
		return 0
	}

	response := models.PrimeEvaluation{
		ParentID:      expression.ParentID,
		Id:            expression.Id,
		Arg1:          expression.Arg1,
		Arg2:          expression.Arg2,
		Operation:     expression.Operation,
		Result:        result,
		OperationTime: operationTime,
		Error:         hasErr,
	}

	p, err := json.Marshal(response)
	if err != nil {
		log.Println("invalid expression:", err)
		return 0
	}

	client := http.Client{}
	_, err = client.Post(fmt.Sprintf("http://%s:%d/internal/task", eh.config.OrchestratorConf.Host, eh.config.OrchestratorConf.Port), //localhost:/internal/task",
		"application/json", bytes.NewReader(p))
	if err != nil {
		log.Println("can't send request to orchestrator")
	} else {
		log.Printf("evaluated expression was successfully send  - %v | result - %f", expression, result)
		return result
	}
	return result

}
