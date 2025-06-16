package profiles

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id,omitempty"`
	DisplayName string             `bson:"display_name,omitempty"`
	Description string             `bson:"description,omitempty"`
	AvatarURL   string             `bson:"avatar_url,omitempty"`
}
