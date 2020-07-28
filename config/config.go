package config

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"

	"github.com/golobby/config"
	"github.com/golobby/config/feeder"
)

type appConfig struct {
	*config.Config
}

func MongoURI() string {
	return ""
}

func NatsURI() string {
	return os.Getenv("NATS_URI", nats.DefaultURL)
}

type Env uint8

const (
	AppEnvDev = iota + 1
	AppEnvProd
)

func AppEnv() Env {
	e, err := C.GetString("app.env")
	if err != nil {
		log.Fatalln(err)
	}
	if e == "prod" || e == "production" {
		return AppEnvProd
	}
	return AppEnvProd
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
