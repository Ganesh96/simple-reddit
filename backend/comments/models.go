package comments

import (
	"time"

	"github.com/ganesh96/simple-reddit/backend/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CommentsCollection *mongo.Collection = configs.GetCollection("comments")
var CommentsVotingHistoryCollection *mongo.Collection = configs.GetCollection("comments_voting_history")

// Comment struct
type Comment struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	PostID       primitive.ObjectID `bson:"post_id,omitempty"`
	Text         string             `bson:"text,omitempty"`
	CreationDate time.Time          `bson:"creation_date,omitempty"`
	UpdationDate time.Time          `bson:"updation_date,omitempty"`
	UpVotes      int                `bson:"up_votes,omitempty"`
	DownVotes    int                `bson:"down_votes,omitempty"`
	Username     string             `bson:"username,omitempty"`
	Edited       bool               `bson:"edited,omitempty"`
}

// VotingHistory represents a user's voting history for a comment
type VotingHistory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CommentID primitive.ObjectID `bson:"comment_id,omitempty"`
	Username  string             `bson:"username,omitempty"`
	Vote      int                `bson:"vote,omitempty"`
}
