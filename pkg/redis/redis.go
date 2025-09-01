package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/tonrock01/another-world-shop/config"
)

func RedisConnect(cfg config.IRedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Password: cfg.Password(),
		DB:       cfg.Db(),
		Protocol: cfg.Protocol(),
	})

	return client
}
