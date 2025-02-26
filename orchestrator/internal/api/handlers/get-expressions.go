package handlers

import (
	"finalTaskLMS/orchestrator/internal/models"
	"finalTaskLMS/orchestrator/types"
	"github.com/gin-gonic/gin"
)

func GetExpressions(o *types.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.JSON(405, gin.H{
				"error": "this method not allowed",
			})
			return
		}
		var response []models.ExpressionsResponse

		o.Mu.Lock()
		defer o.Mu.Unlock()
		for _, v := range o.Queue {
			response = append(response, models.ExpressionsResponse{
				Id:     v.Id,
				Status: v.Status,
				Result: v.Result,
			})
		}

		c.JSON(200, response)
	}

}
