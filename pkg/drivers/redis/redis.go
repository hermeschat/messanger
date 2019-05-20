package redis

import (
	"fmt"
	"strconv"
	"time"

	rds "github.com/go-redis/redis"
)

const Nil = rds.Nil

// ConnectRedis is used to connect to redis
func ConnectRedis() (*rds.Client, error) {
	dbName, err := strconv.Atoi("0")
	if err != nil {
		return nil, err
	}
	Addr := fmt.Sprintf("%s:%s", "127.0.0.1", "6379")
	client := rds.NewClient(&rds.Options{
		Addr:        Addr,
		Password:    "",     // no password set
		DB:          dbName, // use default DB
		DialTimeout: 15 * time.Second,
	})
	return client, nil

}
