package users

import (
	"github.com/ganesh96/simple-reddit/backend/profiles"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User defines the structure for a user in the database.
type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username,omitempty" json:"username,omitempty" validate:"required"`
	Password string             `bson:"password,omitempty" json:"password,omitempty" validate:"required"`
	Profile  profiles.Profile   `bson:"profile,omitempty" json:"profile,omitempty"`
}
