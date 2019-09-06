package main

import (
	"hermes/cmd"
	"hermes/config"
	"hermes/pkg/drivers/redis"

	"github.com/sirupsen/logrus"
)

var env = map[string]string{
	"AUTH_TOKEN":              "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJuYW1laWQiOiJhY2NvdW50cy1zZXJ2aWNlIiwiaWQiOiI1OTgxYTFlNDFkNDFjODRjYWU5MDRmZDMiLCJ1bmlxdWVfbmFtZSI6ImFjY291bnRzLXNlcnZpY2UiLCJzdWIiOiJhY2NvdW50cy1zZXJ2aWNlIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdC8iLCJyb2xlIjpbInpldXMiLCJyb3N0YW0iXSwiYXVkIjoiYjlkYzcxMmM5NTJiNGFhZmI0ODFhYmVkZTBmZWM0ZDgiLCJleHAiOjk5OTk5OTk5OTksIm5iZiI6MTQ5NzE3ODI0NSwiYXBwIjoiNWE5NTYzZjM4NDllMDY3NzEwNWRmNTI5In0.fyU5e4KXpZilnDcxhKRkbYw0paAX15RNGXpifgWvHbY",
	"DATABASE_TYPE":           "mongo",
	"PORT":                    "9000",
	"MONGO_HOST":              "localhost",
	"MONGO_PORT":              "27017",
	"DATABASE_NAME":           "hermes_rc",
	"MONGO_URI":               "mongodb://0.0.0.0:27017/hermes_rc",
	"MONGO_USERNAME":          "",
	"MONGO_PASSWORD":          "",
	"REDIS_HOST":              "localhost",
	"REDIS_PORT":              "6379",
	"REDIS_DB_INDEX":          "5",
	"APPLICATION_SECRET":      "",
	"APPLICATION_SERVICE_URL": "https://api.paygear.ir/application/v3",
	"CLIENT_SECRET":           "J0RYjUcIZHgm41GyPt4wEWUqKzOPXCQAY7n2/ZkQ7WE=",
	"CLIENT_ID":               "b9dc712c952b4aafb481abede0fec4d8",
	"API_KEY":                 "5aa7e856ae7fbc00016ac5a0ede56b6989e14706a6215f4207a40996",
}

func main() {

	cmd.Launch()
	config.Init(env)
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("panic happend:\n%v", err)
		}
		con, err := redis.ConnectRedis()
		if err != nil {
			logrus.Errorf("error while trying to clear redis cache of subscribed channels : %v", err)
			return
		}
		_, err = con.FlushDB().Result()
		if err != nil {
			logrus.Errorf("error while trying to clear redis cache of subscribed channels : %v", err)
			return
		}
	}()
}
