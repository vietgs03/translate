package main

import (
	"log"
	"os"

	"github.com/vietgs03/translate/backend/internal/config"
	"github.com/vietgs03/translate/backend/internal/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Migration direction (up/down) is required")
	}

	direction := os.Args[1]
	if direction != "up" && direction != "down" {
		log.Fatal("Migration direction must be either 'up' or 'down'")
	}

	if err := database.RunMigrations(&cfg.Database, direction); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Printf("Successfully ran migrations: %s", direction)
} 