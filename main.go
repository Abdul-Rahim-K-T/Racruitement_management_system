package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/database"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	secretKey := os.Getenv("SECRET_KEY")
	fmt.Println("Loaded Secret Key:", secretKey)

	// Initialize the database connection
	database.InitDB()

	router := gin.Default()

	routes.LoadAuthRoutes(router)

	router.Run(":8080")
}
