package main

import "finalTaskLMS/agent"

func main() {
	a := agent.NewAgentServer()
	a.InitializeAgent()
	a.ConfigureRouter()
	a.RunServer()
}
