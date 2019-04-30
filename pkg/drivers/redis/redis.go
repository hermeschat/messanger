package redis

import (
	"fmt"
	"strconv"
	"time"

	rds "github.com/go-redis/redis"
)

// ConnectRedis is used to connect to redis
func ConnectRedis() (*rds.Client, error) {
	dbName, err := strconv.Atoi("config.RedisDBName")
	if err != nil {
		return nil, err
	}
	Addr := fmt.Sprintf("%s:%s", "config.RedisHost", "config.RedisPort")
	client := rds.NewClient(&rds.Options{
		Addr:        Addr,
		Password:    "config.RedisPassword", // no password set
		DB:          dbName,                 // use default DB
		DialTimeout: 15 * time.Second,
	})
	return client, nil

}
