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

var communityCollection *mongo.Collection = configs.GetCollection(configs.DB, "communities")

func CreateCommunity() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var community Community

		if err := c.BindJSON(&community); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		count, err := communityCollection.CountDocuments(ctx, bson.M{"name": community.Name})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			common.RespondWithJSON(c, http.StatusConflict, common.COMMUNITY_ALREADY_EXISTS, gin.H{})
			return
		}

		result, err := communityCollection.InsertOne(ctx, community)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"community": result})
	}
}

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

		if communities == nil {
			communities = []Community{}
		}
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"communities": communities})
	}
}

func GetCommunityByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var community Community
		communityName := c.Param("name")

		err := communityCollection.FindOne(ctx, bson.M{"name": communityName}).Decode(&community)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				common.RespondWithJSON(c, http.StatusNotFound, common.COMMUNITY_NOT_FOUND, gin.H{"error": err.Error()})
				return
			}
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"community": community})
	}
}

func DeleteCommunity() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		communityName := c.Param("name")

		result, err := communityCollection.DeleteOne(ctx, bson.M{"name": communityName})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		if result.DeletedCount == 0 {
			common.RespondWithJSON(c, http.StatusNotFound, common.COMMUNITY_NOT_FOUND, gin.H{})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
