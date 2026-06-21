package votes

import (
	"time"

	"github.com/ganesh96/simple-reddit/backend/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TargetPost    = "post"
	TargetComment = "comment"
)

var VotesCollection *mongo.Collection = configs.GetCollection("votes")

// Vote stores one user's latest vote for one target. Counters stay denormalized on posts/comments for fast feeds.
type Vote struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetID     primitive.ObjectID `json:"target_id" bson:"target_id,omitempty"`
	TargetType   string             `json:"target_type" bson:"target_type,omitempty"`
	Username     string             `json:"username" bson:"username,omitempty"`
	Value        int                `json:"value" bson:"value"`
	CreationDate time.Time          `json:"creation_date" bson:"creation_date,omitempty"`
	UpdationDate time.Time          `json:"updation_date" bson:"updation_date,omitempty"`
}

type VoteRequest struct {
	Vote int `json:"vote" binding:"required"`
}
