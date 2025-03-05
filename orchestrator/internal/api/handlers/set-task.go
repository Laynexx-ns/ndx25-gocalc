package handlers

import (
	"finalTaskLMS/globals"
	"finalTaskLMS/orchestrator/types"
	"github.com/gin-gonic/gin"
)

func GetTasks(o *types.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notEvaluatedExpression globals.PrimeEvaluation
		if len(o.Queue) > 0 {
			for _, v := range o.Queue {
				if v.OperationTime == 0 {
					notEvaluatedExpression = v
				}
			}
		}
		c.JSON(200, notEvaluatedExpression)
	}
}
