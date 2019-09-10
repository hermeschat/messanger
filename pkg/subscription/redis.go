package subscription

import (
	"fmt"
	"strconv"
	"time"

	"github.com/amirrezaask/config"
	"github.com/go-redis/redis"
)

const Nil = redis.Nil

// Redis is used to connect to redis
func Redis() (*redis.Client, error) {
	dbName, err := strconv.Atoi(config.Get("redis_db_index"))
	if err != nil {
		return nil, err
	}
	Addr := fmt.Sprintf("%s:%s", config.Get("redis_host"), config.Get("redis_port"))
	client := redis.NewClient(&redis.Options{
		Addr:        Addr,
		Password:    "",     // no password set
		DB:          dbName, // use default DB
		DialTimeout: 15 * time.Second,
	})
	return client, nil

}
