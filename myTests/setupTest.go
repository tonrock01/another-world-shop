package myTests

import (
	"encoding/json"

	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/modules/servers"
	"github.com/tonrock01/another-world-shop/pkg/database"
	"github.com/tonrock01/another-world-shop/pkg/redis"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := database.DBConnect(cfg.Db())

	redisClient := redis.RedisConnect(cfg.Redis())

	s := servers.NewServer(cfg, db, redisClient)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJSON(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
