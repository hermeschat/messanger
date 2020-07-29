package redis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/hermeschat/engine/config"
	"time"
)

type Redis interface {
	GetFromRedis(key string) (value string)
	SetToRedis(key string, value string, expTime time.Duration) bool
	DeleteFromRedis(key string) bool
}

type _redis struct {
	*redis.Client
}

func NewRedis() (Redis, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.C.GetString("_redis.host"), config.C.GetString("_redis.port")),
		Password: "",                          // no password set
		DB:       config.C.GetInt("redis.db"), // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	r := &_redis{
		redisClient,
	}
	return r, nil
}

func (r *_redis) GetFromRedis(key string) (value string) {
	redisResult, err := r.Get(key).Result()
	if err != nil {
		fmt.Println("Error on GetUserConfirmation on RedisDB.Get")
		fmt.Println(err)
		return ""
	} else {
		return redisResult
	}
}
func (r *_redis) SetToRedis(key string, value string, expTime time.Duration) bool {
	err := r.Set(key, value, time.Second*expTime).Err()
	if err != nil {
		fmt.Println("error on setting key on _redis DB.")
		return false
	} else {
		fmt.Println("successful setting key on _redis DB.")
		return true
	}
}

func (r *_redis) DeleteFromRedis(key string) bool {
	err := r.Del(key).Err()
	if err != nil {
		fmt.Println("error on deleting key from _redis DB.")
		return false
	} else {
		fmt.Println("successful deletion key from _redis DB.")
		return true
	}
}
