package interceptors

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"ndx/pkg/logger"
	"ndx/pkg/utils"
	"strings"
	"time"
)

const agentContextKey = "isAgentRequest"

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok && len(md["x-agent-request"]) > 0 && md["x-agent-request"][0] == "true" {
			return handler(ctx, req)
		}

		md, ok = metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("unauthorized")
		}

		logger.L().Logf(0, "server: %v", info.Server)

		if len(md["authorization"]) < 1 {
			return nil, errors.New("unauthorized")
		}
		token := strings.TrimSpace(strings.TrimPrefix(md["authorization"][0], "Bearer"))
		if _, err := utils.VerifyToken(token); err != nil {
			return nil, errors.New("unauthorized")
		}

		logger.L().Logf(0, "request: %v, time: %s | info: %v", req, time.Now().String(), info)

		return handler(ctx, req)
	}
}
