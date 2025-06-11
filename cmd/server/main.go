package main

import (
	"log"
	"mozho_chat/internal/api"
	"mozho_chat/internal/db"
	"mozho_chat/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	dbConn := db.InitPostgres(cfg)

	r := api.SetupRouter(dbConn)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
