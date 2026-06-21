package posts

import (
	"log"
	"net/http"
	"time"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreatePost creates a new post.
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": "Invalid request body"})
		return
	}

	now := time.Now()
	newPost := Post{
		ID:            primitive.NewObjectID(),
		Title:         req.Title,
		Text:          req.Text,
		Community:     req.Community,
		Username:      c.GetString("username"),
		CreationDate:  now,
		UpdationDate:  now,
		UpVotes:       0,
		DownVotes:     0,
		CommentsCount: 0,
	}

	_, err := PostCollection.InsertOne(c.Request.Context(), newPost)
	if err != nil {
		log.Printf("Error creating post: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to create post"})
		return
	}

	common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"message": "Post created successfully", "post": newPost})
}

// GetAllPosts retrieves a bounded page of posts.
func GetAllPosts(c *gin.Context) {
	page, err := common.ParsePageRequest(c)
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{}
	if community := c.Query("community"); community != "" {
		communityID, err := primitive.ObjectIDFromHex(community)
		if err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid community ID"})
			return
		}
		filter["community"] = communityID
	}
	if page.HasAfter {
		filter["_id"] = bson.M{"$lt": page.AfterID}
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(page.Limit + 1)

	cursor, err := PostCollection.Find(c.Request.Context(), filter, findOptions)
	if err != nil {
		log.Printf("Error finding posts: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var results []Post
	if err = cursor.All(c.Request.Context(), &results); err != nil {
		log.Printf("Error decoding posts: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	posts, pagination := common.ApplyCursorPage(results, page.Limit)
	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"posts": posts, "pagination": pagination})
}

// GetPostById retrieves a single post by its ID.
func GetPostById(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		log.Printf("Error converting postId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid post ID"})
		return
	}

	var post Post
	err = PostCollection.FindOne(c.Request.Context(), bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		log.Printf("Error finding post: %v", err)
		common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{"error": "Post not found"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"post": post})
}

// UpdatePost updates a post.
func UpdatePost(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		log.Printf("Error converting postId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid post ID"})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": "Invalid request body"})
		return
	}

	filter := bson.M{"_id": postID, "username": c.GetString("username")}
	update := bson.M{"$set": bson.M{"title": req.Title, "text": req.Text, "updation_date": time.Now()}}

	result, err := PostCollection.UpdateOne(c.Request.Context(), filter, update)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to update post"})
		return
	}
	if result.MatchedCount == 0 {
		common.RespondWithJSON(c, http.StatusForbidden, common.FORBIDDEN, gin.H{"error": "Post not found or not owned by user"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Post updated successfully"})
}

// DeletePost deletes a post.
func DeletePost(c *gin.Context) {
	postID, err := primitive.ObjectIDFromHex(c.Param("postId"))
	if err != nil {
		log.Printf("Error converting postId to ObjectID: %v", err)
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_PARAM, gin.H{"error": "Invalid post ID"})
		return
	}

	result, err := PostCollection.DeleteOne(c.Request.Context(), bson.M{"_id": postID, "username": c.GetString("username")})
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to delete post"})
		return
	}
	if result.DeletedCount == 0 {
		common.RespondWithJSON(c, http.StatusForbidden, common.FORBIDDEN, gin.H{"error": "Post not found or not owned by user"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Post deleted successfully"})
}
