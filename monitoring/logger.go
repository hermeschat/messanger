package monitoring

import (
	"github.com/hermeschat/engine/config"
	"go.uber.org/zap"
)

var sugaredLoggerInstance *zap.SugaredLogger
var LoggerInstance *zap.Logger

func Logger() *zap.SugaredLogger {
	var err error
	if sugaredLoggerInstance != nil {
		return sugaredLoggerInstance
	}
	if config.AppEnv() == config.AppEnvDev {
		if LoggerInstance == nil {
			LoggerInstance, err = zap.NewDevelopment()
			if err != nil {
				panic(err)
			}
		}
		sugaredLoggerInstance = LoggerInstance.Sugar()
		return sugaredLoggerInstance
	} else {
		if LoggerInstance == nil {
			LoggerInstance, err = zap.NewProduction()
			if err != nil {
				panic(err)
			}
		}
		sugaredLoggerInstance = LoggerInstance.Sugar()
		return sugaredLoggerInstance
	}
}
