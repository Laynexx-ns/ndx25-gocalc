package handlers

import (
	"finalTaskLMS/orchestrator/internal/models"
	"finalTaskLMS/orchestrator/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

func AddExpressionHandler(o *types.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.JSON(405, gin.H{
				"error": "this method not allowed",
			})
		}

		var expression models.UserExpressions
		if err := c.ShouldBindJSON(&expression); err != nil {
			c.JSON(400, gin.H{
				"error": "invalid JSON",
			})
			return
		}

		valid, err := regexp.MatchString("^[0-9)(*/+-]+$", expression.Expression)
		if err != nil || !valid {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "request contains forbidden symbols",
			})
			return
		}

		o.Mu.Lock()
		defer o.Mu.Unlock()
		o.Queue = append(o.Queue, models.Expressions{
			Id:         len(o.Queue),
			Status:     "pending",
			Result:     0,
			Expression: expression.Expression,
		})
		c.JSON(200, gin.H{
			"id": o.Queue[len(o.Queue)-1].Id,
		})
	}
}
