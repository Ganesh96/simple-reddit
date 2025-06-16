package profiles

import (
	"context"
	"net/http"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var profileCollection *mongo.Collection = configs.GetCollection(configs.DB, "profiles")

type Profile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id,omitempty"`
	DisplayName string             `bson:"display_name,omitempty"`
	Description string             `bson:"description,omitempty"`
	AvatarURL   string             `bson:"avatar_url,omitempty"`
}

func GetProfileByUsername(c *gin.Context) {
	username := c.Param("username")

	var user common.User
	userCollection := configs.GetCollection(configs.DB, "users")
	err := userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		common.RespondWithJSON(c, http.StatusNotFound, common.USER_NOT_FOUND, gin.H{"error": "User not found"})
		return
	}

	var profile Profile
	err = profileCollection.FindOne(context.TODO(), bson.M{"user_id": user.ID}).Decode(&profile)
	if err != nil {
		common.RespondWithJSON(c, http.StatusNotFound, common.USER_NOT_FOUND, gin.H{"error": "Profile not found"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"profile": profile})
}
