package posts

import (
	"time"

	"github.com/ganesh96/simple-reddit/backend/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var PostCollection *mongo.Collection = configs.GetCollection("posts")
var SavedCollection *mongo.Collection = configs.GetCollection("saved")

// Post stores the feed item plus denormalized counters needed for fast reads.
type Post struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title         string             `json:"title" bson:"title,omitempty"`
	Text          string             `json:"text" bson:"text,omitempty"`
	Community     primitive.ObjectID `json:"community" bson:"community,omitempty"`
	Username      string             `json:"username" bson:"username,omitempty"`
	UpdationDate  time.Time          `json:"updation_date" bson:"updation_date,omitempty"`
	CreationDate  time.Time          `json:"creation_date" bson:"creation_date,omitempty"`
	UpVotes       int                `json:"up_votes" bson:"up_votes"`
	DownVotes     int                `json:"down_votes" bson:"down_votes"`
	CommentsCount int                `json:"comments_count" bson:"comments_count"`
}

func (p Post) GetID() primitive.ObjectID {
	return p.ID
}

type CreatePostRequest struct {
	Title     string             `json:"title" binding:"required,min=1,max=180"`
	Text      string             `json:"text" binding:"max=10000"`
	Community primitive.ObjectID `json:"community" binding:"required"`
}

type UpdatePostRequest struct {
	Title string `json:"title" binding:"required,min=1,max=180"`
	Text  string `json:"text" binding:"max=10000"`
}

// Saved represents a saved post or comment.
type Saved struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ItemID   primitive.ObjectID `json:"item_id" bson:"item_id,omitempty"`
	ItemType string             `json:"item_type" bson:"item_type,omitempty"`
	Username string             `json:"username" bson:"username,omitempty"`
}
