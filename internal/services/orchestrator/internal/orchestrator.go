package internal

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"ndx/internal/services/orchestrator/internal/types"
	orchestratorservice "ndx/pkg/api/orchestrator-service"
	"ndx/pkg/config"
)

type OrchestratorServer struct {
	orchestratorservice.UnimplementedOrchestratorServiceServer
	Orchestrator *types.Orchestrator

	db     *sql.DB
	config config.Config
}

func NewOrchestratorServer(cfg config.Config, db *sql.DB) *OrchestratorServer {
	return &OrchestratorServer{
		Orchestrator: types.NewOrchestrator(),
		db:           db,
		config:       cfg,
	}
}

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

//func (os *OrchestratorServer) Create
