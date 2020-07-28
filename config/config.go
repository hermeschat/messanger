package config

import (
	"fmt"
	"github.com/golobby/config"
	"github.com/golobby/config/feeder"
	"github.com/hermeschat/engine/monitoring"
	"log"
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
func StanURI() (string, string) {
	clusterId, err := C.GetString("CLUSTER_ID")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating NatsURI: %s", err)
	}
	clientId, err := C.GetString("CLIENT_ID")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating NatsURI: %s", err)
	}
	return clusterId, clientId
}
func PostgresURI() string {
	host, err := C.GetString("POSTGRES_HOST")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating postgres uri: %s", err)
	}
	port, err := C.GetString("POSTGRES_PORT")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating postgres uri: %s", err)
	}
	user, err := C.GetString("POSTGRES_USER")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating postgres uri: %s", err)
	}
	pass, err := C.GetString("POSTGRES_PASS")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating postgres uri: %s", err)
	}
	database, err := C.GetString("POSTGRES_DB")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating postgres uri: %s", err)
	}
	sslmode, err := C.GetString("POSTGRES_SSL")
	if err != nil {
		monitoring.Logger().Fatalf("Error in creating postgres uri: %s", err)
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, database, sslmode)
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
