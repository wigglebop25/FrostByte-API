package main

import (
	"log"
	"os"

	"frostbyte-api/internal/config"
	"frostbyte-api/internal/database"
	"frostbyte-api/internal/seeder"
)

func main() {
	log.Println("=== FrostByte Database Seeder ===")

	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Database
	database.Connect(cfg.DSN())

	// Create seeder and run
	dbSeeder := seeder.NewSeeder(database.DB)
	if err := dbSeeder.SeedIfEmpty(); err != nil {
		log.Printf("Seeding failed: %v", err)
		os.Exit(1)
	}

	log.Println("Seeding completed successfully!")
	os.Exit(0)
}
