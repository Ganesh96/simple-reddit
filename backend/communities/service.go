package communities

import (
	"context"
	"net/http"
	"time"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCommunity(c *gin.Context) {
	var community Community
	if err := c.ShouldBindJSON(&community); err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
		return
	}

	// Check if community with the same name already exists
	count, err := CommunityCollection.CountDocuments(context.TODO(), bson.M{"name": community.Name})
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Error checking for existing community"})
		return
	}
	if count > 0 {
		common.RespondWithJSON(c, http.StatusConflict, common.COMMUNITY_ALREADY_EXISTS, gin.H{"error": "Community with this name already exists"})
		return
	}

	community.ID = primitive.NewObjectID()
	community.CreationDate = time.Now()
	community.UpdationDate = time.Now()

	_, err = CommunityCollection.InsertOne(context.TODO(), community)
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to create community"})
		return
	}

	common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"message": "Community created successfully", "community": community})
}

func GetAllCommunities(c *gin.Context) {
	var communities []Community
	cursor, err := CommunityCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to retrieve communities"})
		return
	}

	if err = cursor.All(context.TODO(), &communities); err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to decode communities"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"communities": communities})
}
