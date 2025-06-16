package comments

import (
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

	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid post ID"})
		return
	}

	comment.ID = primitive.NewObjectID()
	comment.PostID = postID
	comment.CreationDate = time.Now()
	comment.UpdationDate = time.Now()
	comment.Edited = false
	comment.UpVotes = 1
	comment.DownVotes = 0
	comment.Username = c.GetString("username") // Get username from JWT

	_, err = CommentsCollection.InsertOne(c.Request.Context(), comment)
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
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM,
			gin.H{"error": "Invalid post ID"})
		return
	}

	var comments []Comment
	cursor, err := CommentsCollection.Find(c.Request.Context(), bson.M{"post_id": postID})
	if err != nil {
		log.Printf("Error finding comments: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to retrieve comments"})
		return
	}

	if err = cursor.All(c.Request.Context(), &comments); err != nil {
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

	// Check if the user is the author of the comment
	var existingComment Comment
	err = CommentsCollection.FindOne(c.Request.Context(), bson.M{"_id": commentID}).Decode(&existingComment)
	if err != nil {
		common.RespondWithJSON(c, http.StatusNotFound, common.COMMENT_NOT_FOUND, gin.H{"error": "Comment not found"})
		return
	}

	authUsername := c.GetString("username")
	if existingComment.Username != authUsername {
		common.RespondWithJSON(c, http.StatusForbidden, common.FORBIDDEN, gin.H{"error": "You are not authorized to update this comment"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"text":          updatedComment.Text,
			"updation_date": time.Now(),
			"edited":        true,
		},
	}

	_, err = CommentsCollection.UpdateOne(c.Request.Context(), bson.M{"_id": commentID}, update)
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
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM,
			gin.H{"error": "Invalid comment ID"})
		return
	}

	var commentToDelete Comment
	err = CommentsCollection.FindOne(c.Request.Context(), bson.M{"_id": commentID}).Decode(&commentToDelete)
	if err != nil {
		common.RespondWithJSON(c, http.StatusNotFound, common.COMMENT_NOT_FOUND, gin.H{"error": "Comment not found"})
		return
	}

	authUsername := c.GetString("username")
	if commentToDelete.Username != authUsername {
		common.RespondWithJSON(c, http.StatusForbidden, common.FORBIDDEN, gin.H{"error": "You are not authorized to delete this comment"})
		return
	}

	_, err = CommentsCollection.DeleteOne(c.Request.Context(), bson.M{"_id": commentID})
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to delete comment"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Comment deleted successfully"})
}
