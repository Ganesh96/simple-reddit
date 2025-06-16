package profiles

import (
	"context"
	"net/http"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateProfile(c *gin.Context) {
	username := c.Param("username")

	var user common.User
	userCollection := configs.GetCollection(configs.DB, "users")
	err := userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		common.RespondWithJSON(c, http.StatusNotFound, common.USER_NOT_FOUND, gin.H{"error": "User not found"})
		return
	}

	var profile Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"display_name": profile.DisplayName,
			"description":  profile.Description,
			"avatar_url":   profile.AvatarURL,
		},
	}

	_, err = profileCollection.UpdateOne(context.TODO(), bson.M{"user_id": user.ID}, update)
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to update profile"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "Profile updated successfully"})
}
