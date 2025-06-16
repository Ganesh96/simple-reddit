package main

import (
	"log"

	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/ganesh96/simple-reddit/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Set up CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//Run the database
	configs.ConnectDB()

	//Setup routes
	routes.SetupRoutes(router)

	log.Fatal(router.Run(":8080"))
}
