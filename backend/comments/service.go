package comments

import (
	"log"
	"net/http"
	"time"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/posts"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateComment creates a new comment.
func CreateComment(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid post ID"})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": "Invalid request body"})
		return
	}

	now := time.Now()
	comment := Comment{
		ID:           primitive.NewObjectID(),
		PostID:       postID,
		Text:         req.Text,
		CreationDate: now,
		UpdationDate: now,
		Edited:       false,
		UpVotes:      0,
		DownVotes:    0,
		Username:     c.GetString("username"),
	}

	_, err = CommentsCollection.InsertOne(c.Request.Context(), comment)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to create comment"})
		return
	}

	_, _ = posts.PostCollection.UpdateOne(c.Request.Context(), bson.M{"_id": postID}, bson.M{"$inc": bson.M{"comments_count": 1}})

	common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"message": "Comment created successfully", "comment": comment})
}

// GetCommentsByPostId retrieves a bounded page of comments for a post.
func GetCommentsByPostId(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid post ID"})
		return
	}

	page, err := common.ParsePageRequest(c)
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"post_id": postID}
	if page.HasAfter {
		filter["_id"] = bson.M{"$gt": page.AfterID}
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetLimit(page.Limit + 1)

	cursor, err := CommentsCollection.Find(c.Request.Context(), filter, findOptions)
	if err != nil {
		log.Printf("Error finding comments: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to retrieve comments"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var results []Comment
	if err = cursor.All(c.Request.Context(), &results); err != nil {
		log.Printf("Error decoding comments: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	comments, pagination := common.ApplyCursorPage(results, page.Limit)
	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"comments": comments, "pagination": pagination})
}

// UpdateComment updates a comment.
func UpdateComment(c *gin.Context) {
	commentID, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": "Invalid request body"})
		return
	}

	filter := bson.M{"_id": commentID, "username": c.GetString("username")}
	update := bson.M{"$set": bson.M{"text": req.Text, "updation_date": time.Now(), "edited": true}}

	result, err := CommentsCollection.UpdateOne(c.Request.Context(), filter, update)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to update comment"})
		return
	}
	if result.MatchedCount == 0 {
		common.RespondWithJSON(c, http.StatusForbidden, common.FORBIDDEN, gin.H{"error": "Comment not found or not owned by user"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Comment updated successfully"})
}

// DeleteComment deletes a comment.
func DeleteComment(c *gin.Context) {
	commentID, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid comment ID"})
		return
	}

	var existing Comment
	if err := CommentsCollection.FindOne(c.Request.Context(), bson.M{"_id": commentID, "username": c.GetString("username")}).Decode(&existing); err != nil {
		common.RespondWithJSON(c, http.StatusForbidden, common.FORBIDDEN, gin.H{"error": "Comment not found or not owned by user"})
		return
	}

	_, err = CommentsCollection.DeleteOne(c.Request.Context(), bson.M{"_id": commentID, "username": c.GetString("username")})
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to delete comment"})
		return
	}

	_, _ = posts.PostCollection.UpdateOne(c.Request.Context(), bson.M{"_id": existing.PostID}, bson.M{"$inc": bson.M{"comments_count": -1}})

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Comment deleted successfully"})
}
