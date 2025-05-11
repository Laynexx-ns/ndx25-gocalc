package internal

import (
	"database/sql"
	"ndx/internal/models"
	"ndx/internal/services/agent/internal/handlers"
	"ndx/pkg/config"
)

type Agent struct {
	Tasks []models.PrimeEvaluation
}

type AgentServer struct {
	Conf      config.Config
	Agent     *Agent
	db        *sql.DB
	evHandler *handlers.EvaluateHandler
}

func NewAgentServer(c config.Config, db *sql.DB) *AgentServer {
	return &AgentServer{
		Conf: c,
		Agent: &Agent{
			Tasks: []models.PrimeEvaluation{},
		},
		db:        db,
		evHandler: handlers.NewEvaluateHandler(db, c),
	}
}

func Start() {

}
