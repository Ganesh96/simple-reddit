package configs

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is the mongo client instance
var DB *mongo.Client
var once sync.Once

// ConnectDB connects to the MongoDB database. It uses sync.Once to ensure it only runs once.
func ConnectDB() {
	once.Do(func() {
		client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}

		//ping the database
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to MongoDB")
		DB = client
	})
}

// GetCollection returns a collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	ConnectDB() // Ensures the DB is connected before returning the collection
	collection := DB.Database("simple-reddit").Collection(collectionName)
	return collection
}
