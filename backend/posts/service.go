package posts

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

var postCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts")

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var post Post

		if err := c.BindJSON(&post); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		newPost := Post{
			Post_id:   primitive.NewObjectID(),
			Title:     post.Title,
			Text:      post.Text,
			Community: post.Community,
			Username:  post.Username,
		}

		result, err := postCollection.InsertOne(ctx, newPost)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"post": result})
	}
}

func GetAllPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var posts []Post

		cursor, err := postCollection.Find(ctx, bson.M{})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &posts); err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		if posts == nil {
			posts = []Post{}
		}

		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"posts": posts})
	}
}

func GetPostById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var post Post
		postId := c.Param("id")

		objId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_POST_ID, gin.H{"error": err.Error()})
			return
		}

		err = postCollection.FindOne(ctx, bson.M{"post_id": objId}).Decode(&post)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{"error": err.Error()})
				return
			}
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"post": post})
	}
}

func UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var post Post
		postId := c.Param("id")
		objId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_POST_ID, gin.H{"error": err.Error()})
			return
		}

		if err := c.BindJSON(&post); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		update := bson.M{"$set": bson.M{"text": post.Text}}
		result, err := postCollection.UpdateOne(ctx, bson.M{"post_id": objId}, update)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		if result.ModifiedCount == 0 {
			common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{})
			return
		}

		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"modifiedCount": result.ModifiedCount})
	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		postId := c.Param("id")
		objId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_POST_ID, gin.H{"error": err.Error()})
			return
		}

		result, err := postCollection.DeleteOne(ctx, bson.M{"post_id": objId})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		if result.DeletedCount == 0 {
			common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
