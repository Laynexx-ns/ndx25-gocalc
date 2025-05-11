package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log = logger.Sugar()
}

func L() *zap.SugaredLogger {
	if log == nil {
		panic("logger not initialized")
	}
	return log
}
