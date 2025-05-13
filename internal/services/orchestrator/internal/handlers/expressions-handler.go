package handlers

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ndx/internal/models"
	"ndx/internal/services/orchestrator/internal/repo"
	"ndx/internal/services/orchestrator/internal/use_case"
	pb "ndx/pkg/api/orchestrator-service"
	"ndx/pkg/logger"
	"regexp"
	"time"
)

type ExpressionsHandler struct {
	pb.UnimplementedOrchestratorServiceServer
	db       *sql.DB
	exprRepo *repo.ExpressionRepository
	taskRepo *repo.TasksRepository
}

func NewExpressionsHandler(repo *repo.ExpressionRepository, repository *repo.TasksRepository, db *sql.DB) *ExpressionsHandler {
	return &ExpressionsHandler{exprRepo: repo, taskRepo: repository, db: db}
}

func (h *ExpressionsHandler) GetExpressions(ctx context.Context, req *pb.GetExpressionsRequest) (*pb.GetExpressionsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid UUID")
	}

	exprs, err := h.exprRepo.GetExpressions(userID)
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

func (h *ExpressionsHandler) PostExpression(ctx context.Context, req *pb.PostExpressionRequest) (*pb.PostExpressionResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	valid, _ := regexp.MatchString(`^[0-9+\-*/().\s]+$`, req.Expression)
	if !valid {
		return nil, status.Error(codes.InvalidArgument, "invalid expression")
	}

	id, err := h.exprRepo.SaveExpression(models.Expressions{
		UserId:     userID,
		Expression: req.Expression,
		Status:     "pending",
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "db error: "+err.Error())
	}

	use_case.CreateTasks(h.taskRepo)

	return &pb.PostExpressionResponse{Id: int32(id)}, nil
}

func (h *ExpressionsHandler) PostExpressionResult(ctx context.Context,
	req *pb.PostExpressionResultRequest) (*pb.PostExpressionResultResponse, error) {

	query := `UPDATE prime_evaluations 
			set result = $2, operation_time = $3, completed_at = $4
			WHERE id = $1`
	if _, err := h.db.Exec(query, req.Id, req.Result, req.OperationTime, time.Now().String()); err != nil {
		logger.L().Logf(0, "can't update evaluation | err: %v", err)
		return nil, err
	}
	return &pb.PostExpressionResultResponse{}, nil
}

func (h *ExpressionsHandler) GetExpressionById(ctx context.Context, req *pb.GetExpressionByIdRequest) (*pb.GetExpressionByIdResponse, error) {
	expr, err := h.exprRepo.GetExpressionById(int(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "expression not found")
	}

	return &pb.GetExpressionByIdResponse{
		Id:     int32(expr.Id),
		Status: expr.Status,
		Result: float32(expr.Result),
	}, nil
}
