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
	EvHandler *handlers.EvaluateHandler
}

func NewAgentServer(c config.Config, db *sql.DB) *AgentServer {
	a := &types.Agent{
		Tasks: []models.PrimeEvaluation{},
	}
	return &AgentServer{
		Conf:      c,
		Agent:     a,
		db:        db,
		EvHandler: handlers.NewEvaluateHandler(db, c, a),
	}
}
