package internal

import (
	"context"
	"database/sql"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal/types"
	"ndx/internal/services/orchestrator/pkg/calc"
	orchestratorservice "ndx/pkg/api/orchestrator-service"
	"ndx/pkg/config"
	"time"
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

//func (os *OrchestratorServer) Create

// CreateTasks my mega function (I love her)
func CreateTasks(orch *types.Orchestrator) {
	orch.Mu.Lock()
	defer orch.Mu.Unlock()

	for i, v := range orch.Expressions {
		if v.Status != "pending" {
			continue
		}

		orch.Expressions[i].Status = "processing"

		go func(id int, expr models.Expressions) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			resChan := make(chan float64, 1)
			errChan := make(chan error, 1)

			go func() {
				calc.Calc(expr.Expression, resChan, errChan, id, orch)
			}()

			select {
			case res := <-resChan:
				orch.Mu.Lock()
				for j, e := range orch.Expressions {
					if e.Id == id {
						orch.Expressions[j].Status = "successfully calculated"
						orch.Expressions[j].Result = res
						break
					}
				}
				orch.Mu.Unlock()

			case <-errChan:
				orch.Mu.Lock()
				for j, e := range orch.Expressions {
					if e.Id == id {
						orch.Expressions[j].Status = "failed"
						orch.Expressions[j].Result = 0
						break
					}
				}
				orch.Mu.Unlock()

			case <-ctx.Done():
				orch.Mu.Lock()
				for j, e := range orch.Expressions {
					if e.Id == id {
						orch.Expressions[j].Status = "timeout"
						break
					}
				}
				orch.Mu.Unlock()
			}
		}(v.Id, v)
	}
}
