package handlers

import (
	"finalTaskLMS/globals"
	"finalTaskLMS/orchestrator/types"
	"github.com/gin-gonic/gin"
)

func GetTasks(o *types.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notEvaluatedExpression globals.PrimeEvaluation
		o.Queue = o.Queue[1:]
		o.SentEvaluations = append(o.SentEvaluations, notEvaluatedExpression)
		c.JSON(200, notEvaluatedExpression)
	}
}
