package main

import (
	"os"

	"github.com/tonrock01/another-world-shop/config"
	"github.com/tonrock01/another-world-shop/modules/servers"
	"github.com/tonrock01/another-world-shop/pkg/database"
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

	servers.NewServer(cfg, db).Start()
}
