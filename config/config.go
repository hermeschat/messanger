package config

import (
	"fmt"
	"github.com/golobby/config"
	"github.com/hermeschat/engine/monitoring"
	"log"
	"github.com/golobby/config/feeder"
)

type appConfig struct {
	*config.Config
}

func MongoURI() string {
	return ""
}

func NatsURI() string {
	host, err := C.GetString("NATS_HOST")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating NatsURI: %s", err)
	}
	port, err := C.GetString("NATS_PORT")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating NatsURI: %s", err)
	}
	return fmt.Sprintf("%s:%s", host, port)
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
	c, err := config.New(config.Options{
		Feeder: &feeder.Yaml{},
	})

	if err != nil {
		return err
	}
	C = &appConfig{c}
	return nil
}
