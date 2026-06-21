package comments

import (
	"time"

	"github.com/ganesh96/simple-reddit/backend/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CommentsCollection *mongo.Collection = configs.GetCollection("comments")

// Comment stores the discussion item plus denormalized vote counters needed for reads.
type Comment struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PostID       primitive.ObjectID `json:"post_id" bson:"post_id,omitempty"`
	Text         string             `json:"text" bson:"text,omitempty"`
	CreationDate time.Time          `json:"creation_date" bson:"creation_date,omitempty"`
	UpdationDate time.Time          `json:"updation_date" bson:"updation_date,omitempty"`
	UpVotes      int                `json:"up_votes" bson:"up_votes"`
	DownVotes    int                `json:"down_votes" bson:"down_votes"`
	Username     string             `json:"username" bson:"username,omitempty"`
	Edited       bool               `json:"edited" bson:"edited"`
}

func (c Comment) GetID() primitive.ObjectID {
	return c.ID
}

type CreateCommentRequest struct {
	Text string `json:"text" binding:"required,min=1,max=10000"`
}

type UpdateCommentRequest struct {
	Text string `json:"text" binding:"required,min=1,max=10000"`
}
