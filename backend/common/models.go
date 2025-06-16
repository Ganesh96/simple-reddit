package common

import "go.mongodb.org/mongo-driver/bson/primitive"

// User struct represents a user in the database
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
}
