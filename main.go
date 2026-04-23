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
	// Load .env if present (ignored in production where env vars are set directly)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config.ConnectDB()

	if err := config.DB.AutoMigrate(&models.Profile{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	if err := seed.SeedProfiles("seed/profiles.json"); err != nil {
		log.Printf("Seeding warning: %v", err)
	} else {
		log.Println("Database seeded successfully")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	routes.SetupRoutes(r)

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
