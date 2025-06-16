package communities

import (
	"context"
	"net/http"
	"time"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// communityCollection is a package-level variable to interact with the "communities" collection in MongoDB.
var communityCollection *mongo.Collection = configs.GetCollection(configs.DB, "communities")

// CreateCommunity handles the creation of a new community.
func CreateCommunity() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var community Community

		// Bind the incoming JSON to the Community struct.
		if err := c.BindJSON(&community); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		// Check if a community with the same name already exists.
		count, err := communityCollection.CountDocuments(ctx, bson.M{"name": community.Name})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			common.RespondWithJSON(c, http.StatusConflict, common.COMMUNITY_ALREADY_EXISTS, gin.H{})
			return
		}

		// Insert the new community.
		result, err := communityCollection.InsertOne(ctx, community)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// Return a 201 Created status with the new community data.
		common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"community": result})
	}
}

// GetAllCommunities retrieves all communities from the database.
func GetAllCommunities() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var communities []Community

		cursor, err := communityCollection.Find(ctx, bson.M{})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &communities); err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// If no communities found, return an empty array.
		if communities == nil {
			communities = []Community{}
		}

		// Return a 200 OK status with the list of communities.
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"communities": communities})
	}
}

// GetCommunityByName retrieves a single community by its name.
func GetCommunityByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var community Community
		communityName := c.Param("name")

		// Find the community by its name.
		err := communityCollection.FindOne(ctx, bson.M{"name": communityName}).Decode(&community)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				common.RespondWithJSON(c, http.StatusNotFound, common.COMMUNITY_NOT_FOUND, gin.H{"error": err.Error()})
				return
			}
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// Return a 200 OK status with the community data.
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"community": community})
	}
}

// DeleteCommunity deletes a community by its name.
func DeleteCommunity() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		communityName := c.Param("name")

		// Delete the community from the collection.
		result, err := communityCollection.DeleteOne(ctx, bson.M{"name": communityName})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// If no document was deleted, the community was not found.
		if result.DeletedCount == 0 {
			common.RespondWithJSON(c, http.StatusNotFound, common.COMMUNITY_NOT_FOUND, gin.H{})
			return
		}

		// Return 204 No Content on successful deletion.
		c.Status(http.StatusNoContent)
	}
}
