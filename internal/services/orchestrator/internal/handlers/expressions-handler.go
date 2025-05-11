package handlers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal"
	"ndx/internal/services/orchestrator/internal/types"
	"regexp"
	"sync/atomic"

	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ExpressionsHandler struct {
	db   *sql.DB
	orch *types.Orchestrator
}

func NewExpressionsHandler(db *sql.DB, orch *types.Orchestrator) *ExpressionsHandler {
	return &ExpressionsHandler{
		db:   db,
		orch: orch,
	}
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

func (eh *ExpressionsHandler) AddExpression() {
	var expression models.UserExpressions

	valid, _ := regexp.MatchString("^[0-9)(*/+-]+$", expression.Expression)
	if !valid {
		c.JSON(422, gin.H{"error": "invalid characters"})
		return
	}

	eh.orch.Mu.Lock()
	newID := atomic.AddUint64(&s.expressionCounter, 1)
	newExpr := models.Expressions{
		Id:         int(newID),
		Status:     "pending",
		Result:     0,
		Expression: expression.Expression,
	}
	eh.orch.Expressions = append(eh.orch.Expressions, newExpr)
	eh.orch.Mu.Unlock()

	internal.CreateTasks(eh.orch)
}
