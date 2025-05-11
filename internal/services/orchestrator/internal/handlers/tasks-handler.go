package handlers

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ndx/internal/services/orchestrator/internal/repository"
	"ndx/internal/services/orchestrator/internal/types"
	pb "ndx/pkg/api/orchestrator-service"
)

type TasksHandler struct {
	db   *sql.DB
	orch *types.Orchestrator
	repo *repository.TasksRepository
}

func NewTasksHandler(db *sql.DB, orch *types.Orchestrator) *TasksHandler {
	return &TasksHandler{
		db:   db,
		orch: orch,
		repo: repository.NewTaskRepository(db),
	}
}

func (th *TasksHandler) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	notEvaluatedExpression := th.repo.GetTask()
	if notEvaluatedExpression == nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &pb.GetTasksResponse{
		ParentID:      int32(notEvaluatedExpression.ParentID),
		Id:            int64(notEvaluatedExpression.Id),
		Arg1:          float32(notEvaluatedExpression.Arg1),
		Arg2:          float32(notEvaluatedExpression.Arg2),
		Operation:     notEvaluatedExpression.Operation,
		OperationTime: 0,
		Result:        0,
		Error:         false,
		CompletedAt:   notEvaluatedExpression.CompletedAt.String(),
	}, nil
}
