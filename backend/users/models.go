package users

import (
	"github.com/ganesh96/simple-reddit/backend/configs"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
