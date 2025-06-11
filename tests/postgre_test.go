package tests

import (
	"fmt"
	"os"
	"testing"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("failed to load .env file")
	}
}

func TestPostgresConnection(t *testing.T) {
	loadEnv()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get DB object: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("ping failed: %v", err)
	}

	t.Log("Connected to PostgreSQL successfully")
}
