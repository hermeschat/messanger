package monitoring

import (
	"github.com/hermeschat/engine/config"
	"go.uber.org/zap"
)

var loggerInstance *zap.SugaredLogger

func Logger() *zap.SugaredLogger {
	if loggerInstance != nil {
		return loggerInstance
	}
	if config.AppEnv() == config.AppEnvDev {
		l, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		loggerInstance = l.Sugar()
		return loggerInstance
	} else {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		loggerInstance = l.Sugar()
		return loggerInstance
	}
}
