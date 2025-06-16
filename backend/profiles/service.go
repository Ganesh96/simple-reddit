package profiles

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

var profileCollection *mongo.Collection = configs.GetCollection(configs.DB, "profiles")

func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		username := c.Param("username")
		var profile Profile

		err := profileCollection.FindOne(ctx, bson.M{"username": username}).Decode(&profile)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				common.RespondWithJSON(c, http.StatusNotFound, common.USER_NOT_FOUND, gin.H{"error": err.Error()})
				return
			}
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"profile": profile})
	}
}
