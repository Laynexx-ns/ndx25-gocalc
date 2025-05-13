package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"math"
	"ndx/internal/models"
	"ndx/internal/services/agent/internal/types"
	pb "ndx/pkg/api/orchestrator-service"

	"ndx/pkg/config"

	"log"
	"time"
)

const agentContextKey = "isAgentRequest"

type EvaluateHandler struct {
	db     *sql.DB
	config config.Config
	Agent  *types.Agent
	client pb.OrchestratorServiceClient
}

func NewEvaluateHandler(db *sql.DB, c config.Config, agent *types.Agent) *EvaluateHandler {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", c.OrchestratorConf.Host, c.OrchestratorConf.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("could not connect to orchestrator gRPC: %v", err)
	}

	client := pb.NewOrchestratorServiceClient(conn)

	return &EvaluateHandler{
		db:     db,
		config: c,
		Agent:  agent,
		client: client,
	}
}

func (eh *EvaluateHandler) CycleTask() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("x-agent-request", "true"))
			task, err := eh.getTask(ctx)
			eh.Agent.Tasks = append(eh.Agent.Tasks, task)
			if err != nil {
				log.Println("Error getting task:", err)
				continue
			}
			eh.processTask(ctx, &task)

		}
	}
}
func (eh *EvaluateHandler) getTask(ctx context.Context) (models.PrimeEvaluation, error) {
	//ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	ctx = context.WithValue(ctx, agentContextKey, true)

	resp, err := eh.client.GetTasks(ctx, &pb.GetTasksRequest{})
	if err != nil {
		return models.PrimeEvaluation{}, err
	}

	return models.PrimeEvaluation{
		ParentID:  int(resp.ParentID),
		Id:        int(resp.Id),
		Arg1:      float64(resp.Arg1),
		Arg2:      float64(resp.Arg2),
		Operation: resp.Operation,
	}, nil
}

func (eh *EvaluateHandler) processTask(ctx context.Context, expression *models.PrimeEvaluation) float64 {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-agent-request", "true"))

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
			result = expression.Arg1 / expression.Arg2
			operationTime = eh.config.AgentConf.TimeConf.DivisionTime
		} else {
			log.Println("Division by zero")
			hasErr = true
		}
	default:
		log.Println("Unknown operation", expression)
		hasErr = true
	}

	_, err := eh.client.PostExpressionResult(ctx, &pb.PostExpressionResultRequest{
		ParentID:      int32(expression.ParentID),
		Id:            int64(expression.Id),
		Arg1:          float32(expression.Arg1),
		Arg2:          float32(expression.Arg2),
		Operation:     expression.Operation,
		Result:        float32(result),
		OperationTime: int32(operationTime),
		Error:         hasErr,
	})
	if err != nil {
		log.Printf("failed to send result to orchestrator: %v", err)
	} else {
		log.Printf("evaluated expression sent - %v | result - %f", expression, result)
	}

	return result
}
