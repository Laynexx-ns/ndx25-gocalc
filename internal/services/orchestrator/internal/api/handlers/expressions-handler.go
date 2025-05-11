package handlers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/types"

	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ExpressionsHandler struct {
	db   *sql.DB
	orch *types.Orchestrator
}

func (eh *ExpressionsHandler) GetExpressions(c echo.Context) error {
	var response []models.ExpressionsResponse

	eh.orch.Mu.Lock()
	defer eh.orch.Mu.Unlock()
	for _, v := range eh.orch.Expressions {
		response = append(response, models.ExpressionsResponse{
			Id:     v.Id,
			Status: v.Status,
			Result: v.Result,
		})
	}
	return c.JSON(200, response)

}

func (eh *ExpressionsHandler) GetExpressionsById(o *types.Orchestrator) gin.HandlerFunc {
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

func (eh *ExpressionsHandler) GetExpressionResult(o *types.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var result models.PrimeEvaluation
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

func (eh *ExpressionsHandler) GetTasks(o *types.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notEvaluatedExpression models.PrimeEvaluation
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
