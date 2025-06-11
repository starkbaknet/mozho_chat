package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mozho_chat/internal/config"
)

var DB *gorm.DB

func InitPostgres(cfg *config.Config) *gorm.DB {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	log.Println("Postgres connection established.")
	return DB
}
