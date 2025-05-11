package main

import (
	"finalTaskLMS/internal/services/agent"
	"finalTaskLMS/pkg/logger"
)

func main() {
	logger.Init()
	a := agent.NewAgentServer()
	a.InitializeAgent()
	a.ConfigureRouter()
	a.RunServer()
}
