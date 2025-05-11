package handlers

import (
	"database/sql"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal/types"
)

type TasksHandler struct {
	db   *sql.DB
	orch *types.Orchestrator
}

func NewTasksHandler(db *sql.DB, orch *types.Orchestrator) *TasksHandler {
	return &TasksHandler{
		db:   db,
		orch: orch,
	}
}

func (th *TasksHandler) GetTasks() {
	var notEvaluatedExpression models.PrimeEvaluation
	if len(th.orch.Queue) > 0 {
		for _, v := range th.orch.Queue {
			if v.OperationTime == 0 {
				notEvaluatedExpression = v
			}
		}
	}

}
