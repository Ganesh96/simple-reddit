package main

import (
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/ganesh96/simple-reddit/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Connect to the database
	configs.ConnectDB()

	// Setup routes
	routes.SetupRouter(router)

	router.Run("localhost:8080")
}
