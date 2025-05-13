package handlers

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ndx/internal/services/orchestrator/internal/repo"
	pb "ndx/pkg/api/orchestrator-service"
)

type TasksHandler struct {
	db   *sql.DB
	repo *repo.TasksRepository
}

func NewTasksHandler(db *sql.DB) *TasksHandler {
	return &TasksHandler{
		db:   db,
		repo: repo.NewTaskRepository(db),
	}
}

func (th *TasksHandler) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	r, err := th.repo.GetPendingTasks()
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &pb.GetTasksResponse{
		ParentID:      int32(r.ParentID),
		Id:            int64(r.Id),
		Arg1:          float32(r.Arg1),
		Arg2:          float32(r.Arg2),
		Operation:     r.Operation,
		OperationTime: 0,
		Result:        0,
		Error:         false,
		CompletedAt:   "",
	}, nil
}
