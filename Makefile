generate_api:
	protoc \
  	  --proto_path=./pkg/proto \
  	  --go_out=. \
  	  --go-grpc_out=. \
  	  --grpc-gateway_out=. \
  	  pkg/proto/orchestrator-service/orchestrator.proto

	protoc \
	  --proto_path=./pkg/proto \
	  --go_out=. \
	  --go-grpc_out=. \
	  --grpc-gateway_out=. \
	  pkg/proto/agent-service/agent.proto

	protoc \
	  --proto_path=./pkg/proto \
	  --go_out=. \
	  --go-grpc_out=. \
	  pkg/proto/user-service/user.proto








run: generate_api
