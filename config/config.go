package config

import (
	"github.com/golobby/config"
	"github.com/golobby/config/feeder"
)

type appConfig *config.Config
var instance appConfig

func Init() error {
	gc, err := config.New(config.Options{Feeder: &feeder.Yaml{Path: "hermes.yml"}})
	if err != nil {
		return err
	}
	instance = gc
	return nil
}