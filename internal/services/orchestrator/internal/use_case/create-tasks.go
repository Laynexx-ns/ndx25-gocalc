package use_case

import (
	"context"
	"log"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal/repo"
	"ndx/internal/services/orchestrator/pkg/calc"
	"time"
)

func createTasks(repo *repo.TasksRepository) {
	expressions, err := repo.GetPendingTasks()
	if err != nil {
		log.Println("can't fetch tasks:", err)
		return
	}

	for _, expr := range expressions {
		go func(expr models.Expressions) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			resChan := make(chan float64, 1)
			errChan := make(chan error, 1)

			go calc.Calc(expr.Expression, resChan, errChan, expr.Id, repo)

			select {
			case res := <-resChan:
				_ = repo.UpdateExpressionResult(expr.Id, "successfully calculated", res)
			case err := <-errChan:
				_ = repo.UpdateExpressionResult(expr.Id, "failed", 0)
				log.Println("calc error:", err)
			case <-ctx.Done():
				_ = repo.UpdateExpressionStatus(expr.Id, "timeout")
			}
		}(expr)
	}
}
