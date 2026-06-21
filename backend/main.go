package main

import (
	"log"
	"time"

	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/ganesh96/simple-reddit/backend/middleware"
	"github.com/ganesh96/simple-reddit/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.BodySizeLimit(1 << 20))
	router.Use(middleware.RateLimit(120, time.Minute))

	// Set up CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	configs.ConnectDB()
	configs.EnsureIndexes()

	routes.SetupRoutes(router)

	log.Fatal(router.Run(":8080"))
}
