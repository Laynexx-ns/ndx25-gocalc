package handlers

import (
	"finalTaskLMS/globals"
	"finalTaskLMS/orchestrator/types"
	"fmt"
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
		for i, v := range o.Queue {
			if v.Id == result.Id && v.ParentID == result.ParentID {
				o.Queue[i].Result = result.Result
				o.Queue[i].OperationTime = 1
				break
			}
		}

		fmt.Println(o.Queue)
		c.JSON(http.StatusOK, gin.H{})

	}
}
