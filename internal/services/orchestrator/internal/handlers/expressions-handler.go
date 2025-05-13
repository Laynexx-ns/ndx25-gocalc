package handlers

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal"
	"ndx/internal/services/orchestrator/internal/repo"
	"ndx/internal/services/orchestrator/internal/types"
	pb "ndx/pkg/api/orchestrator-service"
	"regexp"
	"sync/atomic"
)

type ExpressionsGRPCHandler struct {
	pb.UnimplementedOrchestratorServiceServer
	repo *repo.ExpressionRepository
	orch *types.Orchestrator
}

func NewExpressionsGRPCHandler(repo *repo.ExpressionRepository, orch *types.Orchestrator) *ExpressionsGRPCHandler {
	return &ExpressionsGRPCHandler{repo: repo, orch: orch}
}

func (h *ExpressionsGRPCHandler) GetExpressions(ctx context.Context, req *pb.GetExpressionsRequest) (*pb.GetExpressionsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid UUID")
	}

	exprs, err := h.repo.GetExpressions(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "db error: "+err.Error())
	}

	var resp pb.GetExpressionsResponse
	for _, e := range exprs {
		resp.Response = append(resp.Response, &pb.ExpressionsResponse{
			Id:         int32(e.Id),
			Status:     e.Status,
			Result:     float32(e.Result),
			Expression: e.Expression,
			UserId:     e.UserId.String(),
		})
	}
	return &resp, nil
}

func (h *ExpressionsGRPCHandler) PostExpression(ctx context.Context, req *pb.PostExpressionRequest) (*pb.PostExpressionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	valid, _ := regexp.MatchString(`^[0-9+\-*/().\s]+$`, req.Expression)
	if !valid {
		return nil, status.Error(codes.InvalidArgument, "invalid expression")
	}

	id, err := h.repo.SaveExpression(models.Expressions{
		UserId:     userID,
		Expression: req.Expression,
		Status:     "pending",
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "db error: "+err.Error())
	}

	newID := atomic.AddUint64(&h.orch.ExpressionCounter, 1)
	h.orch.Mu.Lock()
	h.orch.Expressions = append(h.orch.Expressions, models.Expressions{
		Id:         int(id),
		Status:     "pending",
		Result:     0,
		Expression: req.Expression,
		UserId:     userID,
	})
	h.orch.Mu.Unlock()

	go internal.CreateTasks(h.orch)

	return &pb.PostExpressionResponse{Id: int32(id)}, nil
}

func (h *ExpressionsGRPCHandler) GetExpressionById(ctx context.Context, req *pb.GetExpressionByIdRequest) (*pb.GetExpressionByIdResponse, error) {
	expr, err := h.repo.GetExpressionById(int(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "expression not found")
	}

	return &pb.GetExpressionByIdResponse{
		Id:     int32(expr.Id),
		Status: expr.Status,
		Result: float32(expr.Result),
	}, nil
}
