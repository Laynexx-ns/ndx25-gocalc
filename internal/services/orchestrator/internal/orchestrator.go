package internal

import (
	"context"
	"database/sql"
	"ndx/internal/services/orchestrator/internal/handlers"
	"ndx/internal/services/orchestrator/internal/repo"
	"ndx/internal/services/orchestrator/internal/types"
	pb "ndx/pkg/api/orchestrator-service"
	"ndx/pkg/config"
)

type OrchestratorServer struct {
	pb.UnimplementedOrchestratorServiceServer
	Orchestrator       *types.Orchestrator
	db                 *sql.DB
	config             config.Config
	expressionsHandler *handlers.ExpressionsHandler
	taskHandler        *handlers.TasksHandler
}

func NewOrchestratorServer(cfg config.Config, db *sql.DB) *OrchestratorServer {
	return &OrchestratorServer{
		Orchestrator:       types.NewOrchestrator(),
		db:                 db,
		config:             cfg,
		expressionsHandler: handlers.NewExpressionsHandler(repo.NewExpressionRepository(db), repo.NewTaskRepository(db), db),
		taskHandler:        handlers.NewTasksHandler(db),
	}
}

func (os *OrchestratorServer) GetExpressions(ctx context.Context, req *pb.GetExpressionsRequest) (*pb.GetExpressionsResponse, error) {
	return os.expressionsHandler.GetExpressions(ctx, req)
}

func (os *OrchestratorServer) PostExpression(c context.Context, req *pb.PostExpressionRequest) (*pb.PostExpressionResponse, error) {
	return os.expressionsHandler.PostExpression(c, req)
}

func (os *OrchestratorServer) GetExpressionById(c context.Context, req *pb.GetExpressionByIdRequest) (*pb.GetExpressionByIdResponse, error) {
	return os.expressionsHandler.GetExpressionById(c, req)
}

func (os *OrchestratorServer) GetTasks(c context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	return os.taskHandler.GetTasks(c, req)
}

func (os *OrchestratorServer) PostExpressionResult(c context.Context, req *pb.PostExpressionResultRequest) (*pb.PostExpressionResultResponse, error) {
	return os.expressionsHandler.PostExpressionResult(c, req)
}
