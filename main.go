package main

import (
	"log"
	"os"

	"insight-api/config"
	"insight-api/models"
	"insight-api/routes"
	"insight-api/seed"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignored if not present, e.g. in production)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config.ConnectDB()

	config.DB.AutoMigrate(&models.Profile{})

	// Seed the database with profiles (skips duplicates on re-run)
	if err := seed.SeedProfiles("seed/profiles.json"); err != nil {
		log.Printf("Warning: seeding failed: %v", err)
	} else {
		log.Println("Database seeding complete")
	}

	r := gin.Default()
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
