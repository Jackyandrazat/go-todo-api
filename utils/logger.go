package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger(appEnv string) {
	if appEnv == "production" {
		Logger, _ = zap.NewProduction()
	} else {
		Logger, _ = zap.NewDevelopment()
	}
}
