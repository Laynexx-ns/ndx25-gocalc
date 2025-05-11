package handlers

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal"
	"ndx/internal/services/orchestrator/internal/repository"

	"ndx/internal/services/orchestrator/internal/types"
	pb "ndx/pkg/api/orchestrator-service"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"strconv"
)

type ExpressionsHandler struct {
	db   *sql.DB
	orch *types.Orchestrator
	repo *repository.ExpressionRepository
}

func NewExpressionsHandler(db *sql.DB, orch *types.Orchestrator) *ExpressionsHandler {
	return &ExpressionsHandler{
		db:   db,
		orch: orch,
		repo: repository.NewExpressionRepository(db),
	}
}

func (eh *ExpressionsHandler) GetExpressions(ctx context.Context, req *pb.GetExpressionsRequest) (*pb.GetExpressionsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	headers := md["authorization"]
	if len(headers) == 0 {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	parts := strings.SplitN(headers[0], " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	//TODO: add authentication checking logic to get userId

	_, _ = eh.repo.GetExpressions(uuid.New())

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
