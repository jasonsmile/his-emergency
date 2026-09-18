package main

import (
	"fmt"
	"log"
	"os"

	"emergency-his/server/config"
	"emergency-his/server/database"
	"emergency-his/server/router"
)

func main() {
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "config/config.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.InitMySQL(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := router.SetupRouter(cfg.Server.Mode, db)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server listening on %s (config: %s)", addr, configPath)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
