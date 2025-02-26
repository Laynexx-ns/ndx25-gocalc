package main

import "finalTaskLMS/orchestrator"

func main() {
	os := orchestrator.NewOrchestratorServer()
	os.ConfigureRouter()
	os.RunServer(8080)
}
