package db

import (
	"app/config"
	"fmt"
	"github.com/go-redis/redis/v7"
)

func loadRedisConfigurationFromConfig() (*redis.Options, error) {
	host, err := config.C.GetString("kv.redis.host")
	if err != nil {
		return nil, fmt.Errorf("could'nt create redis instance %w", err)
	}
	port, err := config.C.GetString("kv.redis.port")
	if err != nil {
		return nil, fmt.Errorf("could'nt create redis instance %w", err)
	}
	username, err := config.C.GetString("kv.redis.username")
	if err != nil {
		return nil, fmt.Errorf("could'nt create redis instance %w", err)
	}
	password, err := config.C.GetString("kv.redis.password")
	if err != nil {
		return nil, fmt.Errorf("could'nt create redis instance %w", err)
	}
	database, err := config.C.GetInt("kv.redis.database")
	if err != nil {
		return nil, fmt.Errorf("could'nt create redis instance %w", err)
	}
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: username,
		Password: password,
		DB:       database,
	}, nil
}

var cli *redis.Client

func connect(options *redis.Options) (*redis.Client, error) {
	cli = redis.NewClient(options)
	stat := cli.Ping()
	if stat.Err() != nil {
		return nil, stat.Err()
	}
	return cli, nil
}

func NewRedis() (*redis.Client, error) {
	options, err := loadRedisConfigurationFromConfig()
	if err != nil {
	    return nil, err
	}
	if cli == nil {
		return connect(options)
	}
	if err := cli.Ping().Err(); err != nil {
		return connect(options)
	}
	return cli, nil
}
