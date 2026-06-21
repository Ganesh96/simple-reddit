package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := map[string][]mongo.IndexModel{
		"posts": {
			{Keys: bson.D{{Key: "community", Value: 1}, {Key: "_id", Value: -1}}},
			{Keys: bson.D{{Key: "username", Value: 1}, {Key: "_id", Value: -1}}},
		},
		"comments": {
			{Keys: bson.D{{Key: "post_id", Value: 1}, {Key: "_id", Value: 1}}},
			{Keys: bson.D{{Key: "username", Value: 1}, {Key: "_id", Value: -1}}},
		},
		"votes": {
			{
				Keys:    bson.D{{Key: "target_type", Value: 1}, {Key: "target_id", Value: 1}, {Key: "username", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{Keys: bson.D{{Key: "username", Value: 1}, {Key: "target_type", Value: 1}}},
		},
		"saved": {
			{
				Keys:    bson.D{{Key: "item_type", Value: 1}, {Key: "item_id", Value: 1}, {Key: "username", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		"users": {
			{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
	}

	for collectionName, models := range indexes {
		if _, err := GetCollection(collectionName).Indexes().CreateMany(ctx, models); err != nil {
			log.Printf("failed to create indexes for %s: %v", collectionName, err)
		}
	}
}
