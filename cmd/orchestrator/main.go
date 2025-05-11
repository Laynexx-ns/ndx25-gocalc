package main

import (
	"finalTaskLMS/internal/services/orchestrator"
)

func main() {
	os := orchestrator.NewOrchestratorServer()
	os.ConfigureRouter()
	os.RunServer(8080)
}
