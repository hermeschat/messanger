package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type appConfig struct {
	*viper.Viper
}

func MongoURI() string {
	return ""
}
func NatsUrl() string {
	host := C.GetString("host")
	port := C.GetString("port")
	user := C.GetString("user")
	pass := C.GetString("pass")

	uri := fmt.Sprintf("%s:%s", host, port)
	if user != "" && pass != "" {
		uri = fmt.Sprintf("%s:%s@%s", user, pass, uri)
	}
	return uri
}
func PostgresURI() string {
	host := C.GetString("POSTGRES_HOST")
	port := C.GetString("POSTGRES_PORT")
	user := C.GetString("POSTGRES_USER")
	pass := C.GetString("POSTGRES_PASS")
	database := C.GetString("POSTGRES_DB")
	sslmode := C.GetString("POSTGRES_SSL")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, database, sslmode)
}

type Env uint8

const (
	AppEnvDev = iota + 1
	AppEnvProd
)

func AppEnv() Env {
	e := C.GetString("app.env")
	if e == "prod" || e == "production" {
		return AppEnvProd
	}
	return AppEnvProd
}

var C *appConfig

func Init() error {
	v := viper.New()
	v.SetConfigName("hermes")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	C = &appConfig{v}
	return nil
}
