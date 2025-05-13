package main

import (
	"fmt"
	"google.golang.org/grpc"
	"ndx/internal/services/user-service/internal"
	us "ndx/pkg/api/user-service"
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
		logger.L().Fatalf("can't connect to postgres | err: %v", err)
	}

	svc := internal.NewUserService(pgConn)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.UserServiceConf.Host,
		cfg.UserServiceConf.Port))
	if err != nil {
		logger.L().Fatalf("can't start listening user-service | err: %v", err)
	}

	server := grpc.NewServer()
	us.RegisterUserServiceServer(server, svc)

	if err = server.Serve(lis); err != nil {
		logger.L().Fatalf("can't start serving user-service | err: %v", err)
	}

}
