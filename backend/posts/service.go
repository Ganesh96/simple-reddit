package posts

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

// CreatePost creates a new post
func CreatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		log.Printf("Error binding JSON: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY,
			gin.H{"error": "Invalid request body"})
		return
	}

	newPost := Post{
		ID:            primitive.NewObjectID(),
		Title:         post.Title,
		Text:          post.Text,
		Community:     post.Community,
		Username:      post.Username,
		CreationDate:  time.Now(),
		UpdationDate:  time.Now(),
		UpVotes:       1,
		DownVotes:     0,
		CommentsCount: 0,
	}

	_, err := PostCollection.InsertOne(context.TODO(), newPost)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to create post"})
		return
	}

	common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS,
		gin.H{"message": "Post created successfully", "post": newPost})
}

// GetAllPosts retrieves all posts
func GetAllPosts(c *gin.Context) {
	var posts []Post
	cursor, err := PostCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Error finding posts: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to retrieve posts"})
		return
	}

	if err = cursor.All(context.TODO(), &posts); err != nil {
		log.Printf("Error decoding posts: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to retrieve posts"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"posts": posts})
}

// GetPostById retrieves a single post by its ID
func GetPostById(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		log.Printf("Error converting postId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM,
			gin.H{"error": "Invalid post ID"})
		return
	}

	var post Post
	err = PostCollection.FindOne(context.TODO(), bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		log.Printf("Error finding post: %v", err)
		common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND,
			gin.H{"error": "Post not found"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"post": post})
}

// UpdatePost updates a post
func UpdatePost(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		log.Printf("Error converting postId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM,
			gin.H{"error": "Invalid post ID"})
		return
	}

	var updatedPost Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		log.Printf("Error binding JSON: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY,
			gin.H{"error": "Invalid request body"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":         updatedPost.Title,
			"text":          updatedPost.Text,
			"updation_date": time.Now(),
		},
	}

	_, err = PostCollection.UpdateOne(context.TODO(), bson.M{"_id": postID}, update)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to update post"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Post updated successfully"})
}

// DeletePost deletes a post
func DeletePost(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		log.Printf("Error converting postId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM,
			gin.H{"error": "Invalid post ID"})
		return
	}

	_, err = PostCollection.DeleteOne(context.TODO(), bson.M{"_id": postID})
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR,
			gin.H{"error": "Failed to delete post"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Post deleted successfully"})
}
