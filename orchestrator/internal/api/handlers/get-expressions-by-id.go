package handlers

import (
	"finalTaskLMS/orchestrator/internal/models"
	"finalTaskLMS/orchestrator/types"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetExpressionsById(o *types.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedParamId, _ := strconv.Atoi(c.Param("id"))

		o.Mu.Lock()
		defer o.Mu.Unlock()
		for _, v := range o.Expressions {
			if v.Id == parsedParamId {
				c.JSON(200, models.ExpressionsResponse{
					Id:     v.Id,
					Status: v.Status,
					Result: v.Result,
				})
				return
			}
		}
		c.JSON(404, gin.H{
			"response": "Not found",
		})

	}
}
