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

// postCollection is a package-level variable to interact with the "posts" collection in MongoDB.
var postCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts")

// CreatePost handles the creation of a new post.
// It expects a JSON payload with post details and saves it to the database.
func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set up a context with a timeout to prevent long-running database operations.
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var post Post

		// Bind the incoming JSON payload to the Post struct.
		// If binding fails, return a 400 Bad Request error.
		if err := c.BindJSON(&post); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		// Create a new Post instance with a new unique ObjectID.
		newPost := Post{
			Post_id:   primitive.NewObjectID(),
			Title:     post.Title,
			Text:      post.Text,
			Community: post.Community,
			Username:  post.Username,
			// Other fields like votes, created_at can be initialized here.
		}

		// Insert the new post into the database.
		result, err := postCollection.InsertOne(ctx, newPost)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// On successful creation, return a 201 Created status with the new post's ID.
		common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"post": result})
	}
}

// GetAllPosts retrieves all posts from the database.
func GetAllPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var posts []Post

		// Find all documents in the posts collection.
		cursor, err := postCollection.Find(ctx, bson.M{})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		// Iterate through the cursor and decode each document into a Post struct.
		if err = cursor.All(ctx, &posts); err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// If no posts are found, return an empty array instead of null.
		if posts == nil {
			posts = []Post{}
		}

		// Return a 200 OK status with the list of posts.
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"posts": posts})
	}
}

// GetPostById retrieves a single post by its ID.
func GetPostById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var post Post
		// Get the post ID from the URL parameter.
		postId := c.Param("id")

		// Convert the string ID to a MongoDB ObjectID.
		objId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_POST_ID, gin.H{"error": err.Error()})
			return
		}

		// Find the post by its ID.
		err = postCollection.FindOne(ctx, bson.M{"post_id": objId}).Decode(&post)
		if err != nil {
			// If the post is not found, return a 404 Not Found error.
			if err == mongo.ErrNoDocuments {
				common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{"error": err.Error()})
				return
			}
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// Return a 200 OK status with the post data.
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"post": post})
	}
}

// UpdatePost updates an existing post's text.
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

		// Bind the incoming JSON to the Post struct.
		if err := c.BindJSON(&post); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		// Create an update query to set the new text.
		update := bson.M{"$set": bson.M{"text": post.Text}}

		// Find the post and update it.
		result, err := postCollection.UpdateOne(ctx, bson.M{"post_id": objId}, update)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// If no document was modified, it means the post was not found.
		if result.ModifiedCount == 0 {
			common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{})
			return
		}

		// Return a 200 OK status on successful update.
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"modifiedCount": result.ModifiedCount})
	}
}

// DeletePost deletes a post by its ID.
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

		// Delete the post from the collection.
		result, err := postCollection.DeleteOne(ctx, bson.M{"post_id": objId})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// If no document was deleted, the post was not found.
		if result.DeletedCount == 0 {
			common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{})
			return
		}

		// Return a 204 No Content status on successful deletion.
		// Note: A 204 response should not have a body.
		c.Status(http.StatusNoContent)
	}
}
