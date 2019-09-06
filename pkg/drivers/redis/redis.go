package redis

import (
	"fmt"
	"strconv"
	"time"

	"hermes/config"

	rds "github.com/go-redis/redis"
)

const Nil = rds.Nil

// ConnectRedis is used to connect to redis
func ConnectRedis() (*rds.Client, error) {
	dbName, err := strconv.Atoi(config.C().Get("redis_db_index"))
	if err != nil {
		return nil, err
	}
	Addr := fmt.Sprintf("%s:%s", config.C().Get("redis_host"), config.C().Get("redis_port"))
	client := rds.NewClient(&rds.Options{
		Addr:        Addr,
		Password:    "",     // no password set
		DB:          dbName, // use default DB
		DialTimeout: 15 * time.Second,
	})
	return client, nil

}
