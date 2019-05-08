package main

// pgk/errors
// gorilla/mux

import (
	rds "github.com/go-redis/redis"

	"git.raad.cloud/cloud/hermes/pkg/drivers/redis"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	// nats.Run()
	logrus.Info("Hi. Im hermes")

	// cmd.Launch("harchi")
	conn, err := redis.ConnectRedis()
	if err != nil {
		panic(err)
	}
	status := conn.Set("a", "v", time.Second*1000)
	if status.Err() != nil {
		panic(status.Err())
	}
	res, err := conn.Get("k").Result()
	if err == rds.Nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	logrus.Info(res)
}
