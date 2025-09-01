package main

import (
	"os"

	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/modules/servers"
	"github.com/tonrock01/another-world-shop/pkg/database"
	"github.com/tonrock01/another-world-shop/pkg/redis"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := database.DBConnect(cfg.Db())
	defer db.Close()

	redisClient := redis.RedisConnect(cfg.Redis())
	defer redisClient.Close()

	servers.NewServer(cfg, db, redisClient).Start()
}
