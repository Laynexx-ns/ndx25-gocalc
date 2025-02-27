package handlers

import (
	"finalTaskLMS/globals"
	"finalTaskLMS/orchestrator/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetExpressionResult(o *types.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var result globals.PrimeEvaluation
		if err := c.ShouldBindJSON(&result); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		for _, v := range o.Queue {
			if v.Id == result.Id && v.ParentID == result.ParentID {
				v.Result = result.Result
				v.OperationTime = result.OperationTime
			}
		}

		c.JSON(http.StatusOK, gin.H{})

	}
}
