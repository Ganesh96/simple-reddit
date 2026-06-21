package main

import (
	"log"
	"os"
	"strings"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	configs.ConnectDB()
	configs.EnsureIndexes()

	routes.SetupRoutes(router)

	log.Fatal(router.Run(":" + port()))
}

func port() string {
	if value := os.Getenv("PORT"); value != "" {
		return value
	}
	return "8080"
}

func allowedOrigins() []string {
	value := os.Getenv("ALLOWED_ORIGINS")
	if value == "" {
		return []string{"http://localhost:4200"}
	}

	origins := strings.Split(value, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}
