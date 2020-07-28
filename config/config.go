package config

import (
	"github.com/golobby/config"
	"github.com/golobby/config/feeder"
)

type appConfig struct {
	*config.Config
}

func MongoURI() string {
	return ""
}

var C *appConfig

func Init() error {
	c, err := config.New(config.Options{Feeder: &feeder.Yaml{Path: "hermes.yml"}})
	if err != nil {
		return err
	}
	C = &appConfig{c}
	return nil
}