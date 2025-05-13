package main

import (
	"fmt"
	"google.golang.org/grpc"
	"ndx/internal/services/agent/internal"
	agentservice "ndx/pkg/api/agent-service"
	"ndx/pkg/config"
	postgres "ndx/pkg/db/postrgres"
	"ndx/pkg/logger"
	"net"
)

func main() {
	logger.Init()

	cfg := config.NewConfig()
	pgConn, err := postgres.New(cfg.PgConfig)
	if err != nil {
		logger.L().Fatalf("can't connect to database | err: %v", err)
	}
	srv := internal.NewAgentServer(cfg, pgConn)
	go srv.EvHandler.CycleTask()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.AgentConf.Port))
	if err != nil {
		logger.L().Fatalf("can't start listen agent service port | err: %v", err)
	}

	server := grpc.NewServer()
	agentservice.RegisterAgentServiceServer(server, srv)

	if err = server.Serve(lis); err != nil {
		logger.L().Fatalf("can't start agent server | err: %v", err)
	}
}
