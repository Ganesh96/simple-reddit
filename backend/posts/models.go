package posts

import (
	"time"

	"github.com/ganesh96/simple-reddit/backend/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var PostCollection *mongo.Collection = configs.GetCollection("posts")
var SavedCollection *mongo.Collection = configs.GetCollection("saved")

// Post struct with all fields
type Post struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title,omitempty"`
	Text          string             `bson:"text,omitempty"`
	Community     primitive.ObjectID `bson:"community,omitempty"`
	Username      string             `bson:"username,omitempty"`
	UpdationDate  time.Time          `bson:"updation_date,omitempty"`
	CreationDate  time.Time          `bson:"creation_date,omitempty"`
	UpVotes       int                `bson:"up_votes,omitempty"`
	DownVotes     int                `bson:"down_votes,omitempty"`
	CommentsCount int                `bson:"comments_count,omitempty"`
}

// Saved struct represents a saved item
type Saved struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	ItemID   primitive.ObjectID `bson:"item_id,omitempty"`
	ItemType string             `bson:"item_type,omitempty"` // "post" or "comment"
	Username string             `bson:"username,omitempty"`
}
