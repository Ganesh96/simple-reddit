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

var commentCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments")

// CreateComment handles the creation of a new comment.
func CreateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var comment Comment
		defer cancel()

		if err := c.BindJSON(&comment); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		newComment := Comment{
			Id:        primitive.NewObjectID(),
			Text:      comment.Text,
			PostId:    comment.PostId,
			CreatedBy: comment.CreatedBy,
		}

		result, err := commentCollection.InsertOne(ctx, newComment)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"comment": result})
	}
}

// GetCommentsForPost retrieves all comments associated with a specific post.
func GetCommentsForPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		postId := c.Param("postid")
		var comments []Comment

		cursor, err := commentCollection.Find(ctx, bson.M{"postid": postId})
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
}

// DeleteComment deletes a comment by its ID.
func DeleteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		commentId := c.Param("id")
		objId, err := primitive.ObjectIDFromHex(commentId)
		if err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, "Invalid comment ID", gin.H{"error": err.Error()})
			return
		}

		result, err := commentCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		if result.DeletedCount < 1 {
			common.RespondWithJSON(c, http.StatusNotFound, "Comment with specified ID not found!", gin.H{})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
