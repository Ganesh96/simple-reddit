package main

import (
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/ganesh96/simple-reddit/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Configure CORS to allow requests from the frontend development server.
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Connect to the MongoDB database.
	configs.ConnectDB()

	// Set up the API routes.
	routes.SetupRouter(router)

	// Start the server.
	router.Run("localhost:8080")
}
