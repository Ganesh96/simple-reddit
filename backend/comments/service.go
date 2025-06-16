package comments

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateComment creates a new comment
func CreateComment(c *gin.Context) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Printf("Error binding JSON: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY,
			gin.H{"error": "Invalid request body"})
		return
	}

	comment.ID = primitive.NewObjectID()
	comment.CreationDate = time.Now()
	comment.UpdationDate = time.Now()
	comment.Edited = false
	comment.UpVotes = 1
	comment.DownVotes = 0

	_, err := CommentCollection.InsertOne(context.TODO(), comment)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to create comment"})
		return
	}

	common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS,
		gin.H{"message": "Comment created successfully", "comment": comment})
}

// GetCommentsByPostId retrieves all comments for a given post
func GetCommentsByPostId(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		log.Printf("Error converting postId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM,
			gin.H{"error": "Invalid post ID"})
		return
	}

	var comments []Comment
	cursor, err := CommentCollection.Find(context.TODO(), bson.M{"post_id": postID})
	if err != nil {
		log.Printf("Error finding comments: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to retrieve comments"})
		return
	}

	if err = cursor.All(context.TODO(), &comments); err != nil {
		log.Printf("Error decoding comments: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to retrieve comments"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"comments": comments})
}

// UpdateComment updates a comment
func UpdateComment(c *gin.Context) {
	commentID, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		log.Printf("Error converting commentId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM,
			gin.H{"error": "Invalid comment ID"})
		return
	}

	var updatedComment Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		log.Printf("Error binding JSON: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY,
			gin.H{"error": "Invalid request body"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"text":          updatedComment.Text,
			"updation_date": time.Now(),
			"edited":        true,
		},
	}

	_, err = CommentCollection.UpdateOne(context.TODO(), bson.M{"_id": commentID}, update)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to update comment"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Comment updated successfully"})
}

// DeleteComment deletes a comment
func DeleteComment(c *gin.Context) {
	commentID, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		log.Printf("Error converting commentId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM,
			gin.H{"error": "Invalid comment ID"})
		return
	}

	_, err = CommentCollection.DeleteOne(context.TODO(), bson.M{"_id": commentID})
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to delete comment"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Comment deleted successfully"})
}
