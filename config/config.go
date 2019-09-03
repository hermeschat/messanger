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
	"AUTH_TOKEN":         "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJuYW1laWQiOiJhY2NvdW50cy1zZXJ2aWNlIiwiaWQiOiI1OTgxYTFlNDFkNDFjODRjYWU5MDRmZDMiLCJ1bmlxdWVfbmFtZSI6ImFjY291bnRzLXNlcnZpY2UiLCJzdWIiOiJhY2NvdW50cy1zZXJ2aWNlIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdC8iLCJyb2xlIjpbInpldXMiLCJyb3N0YW0iXSwiYXVkIjoiYjlkYzcxMmM5NTJiNGFhZmI0ODFhYmVkZTBmZWM0ZDgiLCJleHAiOjk5OTk5OTk5OTksIm5iZiI6MTQ5NzE3ODI0NSwiYXBwIjoiNWE5NTYzZjM4NDllMDY3NzEwNWRmNTI5In0.fyU5e4KXpZilnDcxhKRkbYw0paAX15RNGXpifgWvHbY",
	"DATABASE_TYPE":      "mongo",
	"PORT":               "9000",
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
	"CLIENT_SECRET":      "J0RYjUcIZHgm41GyPt4wEWUqKzOPXCQAY7n2/ZkQ7WE=",
	"CLIENT_ID":          "b9dc712c952b4aafb481abede0fec4d8",
	"API_KEY":            "5aa7e856ae7fbc00016ac5a0ede56b6989e14706a6215f4207a40996",
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
