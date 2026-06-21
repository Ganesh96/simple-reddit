package configs

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var loadEnvOnce sync.Once

func loadEnv() {
	loadEnvOnce.Do(func() {
		if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
			log.Printf("error loading .env file: %v", err)
		}
	})
}

func EnvMongoURI() string {
	loadEnv()
	if uri := os.Getenv("MONGOURI"); uri != "" {
		return uri
	}
	return "mongodb://localhost:27017"
}

func SecretKey() string {
	loadEnv()
	return os.Getenv("SECRET_KEY")
}
