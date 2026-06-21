package votes

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ganesh96/simple-reddit/backend/comments"
	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/posts"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func VotePost(c *gin.Context) {
	upsertVote(c, TargetPost, "postId")
}

func DeletePostVote(c *gin.Context) {
	deleteVote(c, TargetPost, "postId")
}

func VoteComment(c *gin.Context) {
	upsertVote(c, TargetComment, "commentId")
}

func DeleteCommentVote(c *gin.Context) {
	deleteVote(c, TargetComment, "commentId")
}

func upsertVote(c *gin.Context, targetType string, paramName string) {
	targetID, err := primitive.ObjectIDFromHex(c.Param(paramName))
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid target ID"})
		return
	}

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil || (req.Vote != 1 && req.Vote != -1) {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": "vote must be 1 or -1"})
		return
	}

	targetCollection, err := targetCollection(targetType)
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid target type"})
		return
	}

	count, err := targetCollection.CountDocuments(c.Request.Context(), bson.M{"_id": targetID})
	if err != nil || count == 0 {
		common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{"error": "Target not found"})
		return
	}

	username := c.GetString("username")
	filter := bson.M{"target_type": targetType, "target_id": targetID, "username": username}

	oldVote := 0
	var existing Vote
	err = VotesCollection.FindOne(c.Request.Context(), filter).Decode(&existing)
	if err == nil {
		oldVote = existing.Value
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("Error finding vote: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to read vote"})
		return
	}

	if oldVote == req.Vote {
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Vote already applied"})
		return
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{"value": req.Vote, "updation_date": now},
		"$setOnInsert": bson.M{
			"_id":           primitive.NewObjectID(),
			"target_type":  targetType,
			"target_id":    targetID,
			"username":     username,
			"creation_date": now,
		},
	}
	_, err = VotesCollection.UpdateOne(c.Request.Context(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error saving vote: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to save vote"})
		return
	}

	inc := voteCounterDelta(oldVote, req.Vote)
	_, err = targetCollection.UpdateOne(c.Request.Context(), bson.M{"_id": targetID}, bson.M{"$inc": inc})
	if err != nil {
		log.Printf("Error updating vote counters: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to update vote counters"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Vote saved successfully", "vote": req.Vote})
}

func deleteVote(c *gin.Context, targetType string, paramName string) {
	targetID, err := primitive.ObjectIDFromHex(c.Param(paramName))
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid target ID"})
		return
	}

	username := c.GetString("username")
	filter := bson.M{"target_type": targetType, "target_id": targetID, "username": username}

	var existing Vote
	if err := VotesCollection.FindOne(c.Request.Context(), filter).Decode(&existing); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Vote already removed"})
			return
		}
		log.Printf("Error finding vote: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to read vote"})
		return
	}

	_, err = VotesCollection.DeleteOne(c.Request.Context(), filter)
	if err != nil {
		log.Printf("Error deleting vote: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to delete vote"})
		return
	}

	targetCollection, err := targetCollection(targetType)
	if err == nil {
		_, _ = targetCollection.UpdateOne(c.Request.Context(), bson.M{"_id": targetID}, bson.M{"$inc": voteCounterDelta(existing.Value, 0)})
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Vote removed successfully"})
}

func voteCounterDelta(oldVote int, newVote int) bson.M {
	inc := bson.M{"up_votes": 0, "down_votes": 0}
	if oldVote == 1 {
		inc["up_votes"] = inc["up_votes"].(int) - 1
	}
	if oldVote == -1 {
		inc["down_votes"] = inc["down_votes"].(int) - 1
	}
	if newVote == 1 {
		inc["up_votes"] = inc["up_votes"].(int) + 1
	}
	if newVote == -1 {
		inc["down_votes"] = inc["down_votes"].(int) + 1
	}
	return inc
}

func targetCollection(targetType string) (*mongo.Collection, error) {
	switch targetType {
	case TargetPost:
		return posts.PostCollection, nil
	case TargetComment:
		return comments.CommentsCollection, nil
	default:
		return nil, errors.New("unsupported target type")
	}
}
