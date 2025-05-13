package use_case

import (
	"context"
	"log"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal/repo"
	"ndx/internal/services/orchestrator/pkg/calc"
	"ndx/pkg/logger"
	"time"
)

func CreateTasks(repo *repo.TasksRepository) {
	expression, err := repo.GetPendingTasks()
	if err != nil {
		logger.L().Logf(0, "can't get pending tasks | err: %v", err)
		return
	}

	for _, expr := range expression {
		go func(expr models.Expressions) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			resChan := make(chan float64, 1)
			errChan := make(chan error, 1)

			go calc.Calc(expr.Expression, resChan, errChan, expr.Id, repo)

			select {
			case res := <-resChan:
				_ = repo.UpdateExpressionResult(expr.Id, "successfully calculated", res)
			case err = <-errChan:
				_ = repo.UpdateExpressionResult(expr.Id, "failed", 0)
				log.Println("calc error:", err)
			case <-ctx.Done():
				_ = repo.UpdateExpressionStatus(expr.Id, "timeout")
			}
		}(expr)
	}
}
