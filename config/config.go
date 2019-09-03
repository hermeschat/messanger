package config

import (
	"os"
	"strings"
)

//ConfigMap represnets config
type ConfigMap map[string]string

func (c ConfigMap) Get(key string) string {
	val, exists := (c)[key]
	if !exists {
		return ""
	}
	return val
}

func (c *ConfigMap) Set(key, value string) {
	(*c)[key] = value
}

//envlist shows whats envs are available
var envList = map[string]string{
	"AUTH_TOKEN":         "",
	"DATABASE_TYPE":      "mongo",
	"DATABASE_HOST":      "localhost",
	"DATABASE_PORT":      "27017",
	"DATABASE_URI":       "mongodb://localhost:27017/hermes",
	"DATABASE_USERNAME":  "",
	"DATABASE_PASSWORD":  "",
	"DATABASE_NAME":      "hermes",
	"REDIS_HOST":         "localhost",
	"REDIS_PORT":         "6379",
	"REDIS_DB_INDEX":     "5",
	"APPLICATION_SECRET": "",
	"CLIENT_SECRET":      "",
}
var config *ConfigMap

//Config gets the global config object
func Config() ConfigMap {
	return *config
}

/*FromEnv gets a map from env key to default value,
tries to get key from env if not found uses default value present as value of map */
func FromEnv(kd map[string]string) error {
	//err := godotenv.Load("config.env", "secrets.env")
	//if err != nil {
	//	return errors.Wrap(err, "error in loading env file")
	//}
	if kd == nil {
		kd = envList
	}
	c := &ConfigMap{}
	for k, d := range kd {
		v := os.Getenv(k)
		if v == "" {
			v = d
		}
		c.Set(strings.ToLower(k), v)
	}
	config = c
	return nil
}
