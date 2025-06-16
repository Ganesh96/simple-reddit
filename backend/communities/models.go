package communities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
    "github.com/ganesh96/simple-reddit/backend/configs"
)

// CommunityCollection is the collection for communities
var CommunityCollection *mongo.Collection = configs.GetCollection(configs.DB, "communities")

// Community struct represents a community in the database
type Community struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name,omitempty"`
	Description   string             `bson:"description,omitempty"`
	CreationDate  time.Time          `bson:"creation_date,omitempty"`
	UpdationDate  time.Time          `bson:"updation_date,omitempty"`
	MembersCount  int                `bson:"members_count,omitempty"`
	PostsCount    int                `bson:"posts_count,omitempty"`
	Creator       primitive.ObjectID `bson:"creator,omitempty"`
}
