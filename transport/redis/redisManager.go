package redis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/hermeschat/engine/config"
	"time"
)

type IRedisManager interface {
	GetFromRedis(key string) (value string)
	SetToRedis(key string, value string, expTime time.Duration) bool
	DeleteFromRedis(key string) bool
}

type MySweetRedis struct {
	*redis.Client
}

func StartRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.C.GetString("redis.host"), config.C.GetString("redis.port")),
		Password: "",                          // no password set
		DB:       config.C.GetInt("redis.db"), // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
	return redisClient
}

func (r *MySweetRedis) GetFromRedis(key string) (value string) {
	redisResult, err := r.Get(key).Result()
	if err != nil {
		fmt.Println("Error on GetUserConfirmation on RedisDB.Get")
		fmt.Println(err)
		return ""
	} else {
		return redisResult
	}
}
func (r *MySweetRedis) SetToRedis(key string, value string, expTime time.Duration) bool {
	err := r.Set(key, value, time.Second*expTime).Err()
	if err != nil {
		fmt.Println("error on setting key on redis DB.")
		return false
	} else {
		fmt.Println("successful setting key on redis DB.")
		return true
	}
}

func (r *MySweetRedis) DeleteFromRedis(key string) bool {
	err := r.Del(key).Err()
	if err != nil {
		fmt.Println("error on deleting key from redis DB.")
		return false
	} else {
		fmt.Println("successful deletion key from redis DB.")
		return true
	}
}
