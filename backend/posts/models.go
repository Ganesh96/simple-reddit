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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define collections
var PostsCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts")
var PostsVotingHistoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts_voting_history")

// Post struct with all fields
type Post struct {
	Post_id        primitive.ObjectID `bson:"post_id,omitempty"`
	Title          string             `bson:"title,omitempty"`
	Text           string             `bson:"text,omitempty"`
	Community      string             `bson:"community,omitempty"`
	Username       string             `bson:"username,omitempty"`
	Updation_date  time.Time          `bson:"updation_date,omitempty"`
	Creation_date  time.Time          `bson:"creation_date,omitempty"`
	Up_votes       int                `bson:"up_votes,omitempty"`
	Down_votes     int                `bson:"down_votes,omitempty"`
	Comments_count int                `bson:"comments_count,omitempty"`
}

// Posts is an array of Post
type Posts []Post

// PostsVotingHistory struct
type PostsVotingHistory struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Post_id    string             `bson:"post_id,omitempty"`
	Username   string             `bson:"username,omitempty"`
	Vote_value int                `bson:"vote_value,omitempty"`
}

// GetPostByID retrieves a post by its ID
func GetPostByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var post Post
	postId := c.Param("id")

	objId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_POST_ID, gin.H{"error": err.Error()})
		return
	}

	err = PostsCollection.FindOne(ctx, bson.M{"post_id": objId}).Decode(&post)
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

// GetAllPosts retrieves all posts
func GetAllPosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var posts []Post
	cursor, err := PostsCollection.Find(ctx, bson.M{})
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

// CreatePost creates a new post
func CreatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var post Post

	if err := c.BindJSON(&post); err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
		return
	}

	newPost := Post{
		Post_id:        primitive.NewObjectID(),
		Title:          post.Title,
		Text:           post.Text,
		Community:      post.Community,
		Username:       post.Username,
		Creation_date:  time.Now(),
		Updation_date:  time.Now(),
		Up_votes:       0,
		Down_votes:     0,
		Comments_count: 0,
	}

	result, err := PostsCollection.InsertOne(ctx, newPost)
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
		return
	}

	common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"post": result})
}

// UpdatePost updates an existing post
func UpdatePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	update := bson.M{
		"$set": bson.M{
			"text":          post.Text,
			"updation_date": time.Now(),
		},
	}

	result, err := PostsCollection.UpdateOne(ctx, bson.M{"post_id": objId}, update)
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount < 1 {
		common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, nil)
		return
	}
	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"modifiedCount": result.ModifiedCount})
}

// DeletePost deletes a post
func DeletePost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	postId := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(postId)

	result, err := PostsCollection.DeleteOne(ctx, bson.M{"post_id": objId})
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount < 1 {
		common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, nil)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdatePostVote updates the vote on a post
func UpdatePostVote(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var votingHistory PostsVotingHistory
	var post Post

	if err := c.BindJSON(&votingHistory); err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
		return
	}

	objId, _ := primitive.ObjectIDFromHex(votingHistory.Post_id)

	// Check if the user has already voted
	err := PostsVotingHistoryCollection.FindOne(ctx, bson.M{"post_id": votingHistory.Post_id, "username": votingHistory.Username}).Decode(&votingHistory)

	if err != nil { // If no record found, create one
		if err == mongo.ErrNoDocuments {
			_, insertErr := PostsVotingHistoryCollection.InsertOne(ctx, votingHistory)
			if insertErr != nil {
				common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": insertErr.Error()})
				return
			}
		}
	} else { // If record found, update it
		update := bson.M{"$set": bson.M{"vote_value": votingHistory.Vote_value}}
		_, updateErr := PostsVotingHistoryCollection.UpdateOne(ctx, bson.M{"post_id": votingHistory.Post_id, "username": votingHistory.Username}, update)
		if updateErr != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": updateErr.Error()})
			return
		}
	}

	// Update the vote count on the post itself
	err = PostsCollection.FindOne(ctx, bson.M{"post_id": objId}).Decode(&post)
	if err != nil {
		common.RespondWithJSON(c, http.StatusNotFound, common.POST_NOT_FOUND, gin.H{"error": err.Error()})
		return
	}

	// Logic to calculate new vote counts
	upVotes, downVotes := calculateVotes(ctx, votingHistory.Post_id)
	update := bson.M{"$set": bson.M{"up_votes": upVotes, "down_votes": downVotes}}
	_, updateErr := PostsCollection.UpdateOne(ctx, bson.M{"post_id": objId}, update)
	if updateErr != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": updateErr.Error()})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"up_votes": upVotes, "down_votes": downVotes})
}

// helper function to calculate votes
func calculateVotes(ctx context.Context, postId string) (upVotes int, downVotes int) {
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$post_id"},
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

	cursor, err := PostsVotingHistoryCollection.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		return 0, 0
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return 0, 0
	}

	for _, result := range results {
		if result["_id"] == postId {
			upVotes = int(result["up_votes"].(int32))
			downVotes = int(result["down_votes"].(int32))
			return
		}
	}
	return 0, 0
}
