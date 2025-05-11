package config

import (
	postgres "ndx/pkg/db/postrgres"
)

type Config struct {
	AgentConf        AgentConf          `yaml:"AGENT"`
	OrchestratorConf OrchestratorConfig `yaml:"ORCHESTRATOR"`
	PgConfig         postgres.Config    `yaml:"POSTGRES"`
}

type AgentConf struct {
	Port     int           `yaml:"AGENT_PORT" yaml-default:"8082"`
	Host     string        `yaml:"AGENT_HOST" yaml-default:"localhost"`
	TimeConf AgentTimeConf `yaml:"AGENT_TIME_CONF"`
}

type AgentTimeConf struct {
	AdditionTime       int `yaml:"TIME_ADDITION_MS" yaml-default:"50"`
	SubtractionTime    int `yaml:"TIME_SUBTRACTION_MS" yaml-default:"50"`
	MultiplicationTime int `yaml:"TIME_MULTIPLICATIONS_MS" yaml-default:"50"`
	DivisionTime       int `yaml:"TIME_DIVISIONS_MS" yaml-default:"50"`
}

type OrchestratorConfig struct {
	Port int    `yaml:"ORCHESTRATOR_PORT" yaml-default:"8081"`
	Host string `yaml:"ORCHESTRATOR_HOST" yaml-default:"localhost"`
}
