package main

import (
	"fmt"
	"google.golang.org/grpc"
	"ndx/internal/services/orchestrator/internal"
	orchestratorservice "ndx/pkg/api/orchestrator-service"
	"ndx/pkg/config"
	postgres "ndx/pkg/db/postrgres"
	"ndx/pkg/logger"
	"net"
)

func main() {
	logger.Init()

	cfg := config.NewConfig()
	logger.L().Logf(0, "config: %v", cfg)

	pgConn, err := postgres.New(cfg.PgConfig)
	if err != nil {
		logger.L().Fatalf("can't connect to postgres | err: %v", err)
	}

	svc := internal.NewOrchestratorServer(cfg, pgConn)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.OrchestratorConf.Port))
	if err != nil {
		logger.L().Fatalf("can't start listening on orch port | err: %v", err)
	}

	server := grpc.NewServer()
	orchestratorservice.RegisterOrchestratorServiceServer(server, svc)

	if err = server.Serve(lis); err != nil {
		logger.L().Fatalf("can't start serving | err: %v", err)
	}

}
