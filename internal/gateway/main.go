package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	os "ndx/pkg/api/orchestrator-service"
	us "ndx/pkg/api/user-service"
	"ndx/pkg/config"
	"ndx/pkg/logger"
	"ndx/pkg/middlewares"
	"net/http"
)

func main() {
	logger.Init()

	cfg := config.NewConfig()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
			middleware.SetCookieFromContext(ctx, w)
			return nil
		}),
	)

	createUserServiceConnection(ctx, mux, cfg)
	createOrchestratorConnection(ctx, mux, cfg)

	logger.L().Logf(0, "gateway started at the port: %d", cfg.GatewayPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GatewayPort), mux); err != nil {
		logger.L().Fatalf("can't start gateway | err: %v", err)
	}

}

func createUserServiceConnection(ctx context.Context, mux *runtime.ServeMux, cfg config.Config) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.UserServiceConf.Host,
		cfg.UserServiceConf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.L().Fatalf("can't connect to user-service | err: %v", err)
	}

	client := us.NewUserServiceClient(conn)
	if err = us.RegisterUserServiceHandlerClient(ctx, mux, client); err != nil {
		logger.L().Fatalf("failed to register user-service client | err: %v", err)
	}
}

func createOrchestratorConnection(ctx context.Context, mux *runtime.ServeMux, cfg config.Config) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.OrchestratorConf.Host,
		cfg.OrchestratorConf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.L().Fatalf("can't connect to orch | err: %v", err)
	}

	client := os.NewOrchestratorServiceClient(conn)
	if err = os.RegisterOrchestratorServiceHandlerClient(ctx, mux, client); err != nil {
		logger.L().Fatalf("failed to register orch handler client | err: %v", err)
	}
}
