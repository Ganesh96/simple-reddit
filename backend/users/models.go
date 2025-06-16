package users

import (
	"github.com/ganesh96/simple-reddit/backend/configs"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection("users")

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
