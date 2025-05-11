package internal

import (
	"database/sql"
	"ndx/internal/models"
	"ndx/internal/services/agent/internal/handlers"
	"ndx/internal/services/agent/internal/types"
	agentservice "ndx/pkg/api/agent-service"
	"ndx/pkg/config"
)

type AgentServer struct {
	agentservice.UnimplementedAgentServiceServer
	Conf      config.Config
	Agent     *types.Agent
	db        *sql.DB
	evHandler *handlers.EvaluateHandler
}

func NewAgentServer(c config.Config, db *sql.DB) *AgentServer {
	return &AgentServer{
		Conf: c,
		Agent: &types.Agent{
			Tasks: []models.PrimeEvaluation{},
		},
		db:        db,
		evHandler: handlers.NewEvaluateHandler(db, c),
	}
}
