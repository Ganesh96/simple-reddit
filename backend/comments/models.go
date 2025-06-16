package comments

import (
	"context"
	"net/http"
	"time"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CommentsCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments")
var CommentsVotingHistoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments_voting_history")

// Comment struct
type Comment struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Post_id       string             `bson:"post_id,omitempty"`
	Text          string             `bson:"text,omitempty"`
	Creation_date time.Time          `bson:"creation_date,omitempty"`
	Updation_date time.Time          `bson:"updation_date,omitempty"`
	Up_votes      int                `bson:"up_votes,omitempty"`
	Down_votes    int                `bson:"down_votes,omitempty"`
	Username      string             `bson:"username,omitempty"`
	Edited        bool               `bson:"edited,omitempty"`
}

type Comments []Comment

type CommentsVotingHistory struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Comment_id string             `bson:"comment_id,omitempty"`
	Username   string             `bson:"username,omitempty"`
	Vote_value int                `bson:"vote_value,omitempty"`
}

// GetCommentsForPost gets all comments for a post
func GetCommentsForPost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	postId := c.Param("postid")
	var comments []Comment

	cursor, err := CommentsCollection.Find(ctx, bson.M{"post_id": postId})
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &comments); err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
		return
	}
	if comments == nil {
		comments = []Comment{}
	}
	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"comments": comments})
}

// UpdateCommentVote updates the vote on a comment
func UpdateCommentVote(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var votingHistory CommentsVotingHistory
	var comment Comment

	if err := c.BindJSON(&votingHistory); err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
		return
	}

	objId, _ := primitive.ObjectIDFromHex(votingHistory.Comment_id)
	err := CommentsVotingHistoryCollection.FindOne(ctx, bson.M{"comment_id": votingHistory.Comment_id, "username": votingHistory.Username}).Decode(&votingHistory)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, insertErr := CommentsVotingHistoryCollection.InsertOne(ctx, votingHistory)
			if insertErr != nil {
				common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": insertErr.Error()})
				return
			}
		}
	} else {
		update := bson.M{"$set": bson.M{"vote_value": votingHistory.Vote_value}}
		_, updateErr := CommentsVotingHistoryCollection.UpdateOne(ctx, bson.M{"comment_id": votingHistory.Comment_id, "username": votingHistory.Username}, update)
		if updateErr != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": updateErr.Error()})
			return
		}
	}

	err = CommentsCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&comment)
	if err != nil {
		common.RespondWithJSON(c, http.StatusNotFound, "Comment not found", gin.H{"error": err.Error()})
		return
	}

	upVotes, downVotes := calculateCommentVotes(ctx, votingHistory.Comment_id)
	update := bson.M{"$set": bson.M{"up_votes": upVotes, "down_votes": downVotes}}
	_, updateErr := CommentsCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if updateErr != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": updateErr.Error()})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"up_votes": upVotes, "down_votes": downVotes})
}

// helper function to calculate comment votes
func calculateCommentVotes(ctx context.Context, commentId string) (upVotes int, downVotes int) {
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$comment_id"},
			{"up_votes", bson.D{
				{"$sum", bson.D{
					{"$cond", bson.A{bson.D{{"$eq", bson.A{"$vote_value", 1}}}, 1, 0}},
				}},
			}},
			{"down_votes", bson.D{
				{"$sum", bson.D{
					{"$cond", bson.A{bson.D{{"$eq", bson.A{"$vote_value", -1}}}, 1, 0}},
				}},
			}},
		}},
	}
	cursor, err := CommentsVotingHistoryCollection.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		return 0, 0
	}
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return 0, 0
	}
	for _, result := range results {
		if result["_id"] == commentId {
			upVotes = int(result["up_votes"].(int32))
			downVotes = int(result["down_votes"].(int32))
			return
		}
	}
	return 0, 0
}
