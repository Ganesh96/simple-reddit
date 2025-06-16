package profiles

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

// Define the collection from the database
var ProfilesCollection *mongo.Collection = configs.GetCollection(configs.DB, "profiles")

// Profile struct
type Profile struct {
	Id                      primitive.ObjectID `bson:"_id,omitempty"`
	Username                string             `bson:"username,omitempty"`
	User_karma              int                `bson:"user_karma,omitempty"`
	Bio                     string             `bson:"bio,omitempty"`
	Cake_day                time.Time          `bson:"cake_day,omitempty"`
	Number_of_posts         int                `bson:"number_of_posts,omitempty"`
	Number_of_comments      int                `bson:"number_of_comments,omitempty"`
	Number_of_saved_posts   int                `bson:"number_of_saved_posts,omitempty"`
	Number_of_saved_cmnts   int                `bson:"number_of_saved_cmnts,omitempty"`
	Number_of_upvoted_posts int                `bson:"number_of_upvoted_posts,omitempty"`
	Number_of_dnvoted_posts int                `bson:"number_of_dnvoted_posts,omitempty"`
	Number_of_upvoted_cmnts int                `bson:"number_of_upvoted_cmnts,omitempty"`
	Number_of_dnvoted_cmnts int                `bson:"number_of_dnvoted_cmnts,omitempty"`
}

// GetProfile retrieves a user's profile
func GetProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.Param("username")
	var profile Profile

	err := ProfilesCollection.FindOne(ctx, bson.M{"username": username}).Decode(&profile)
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
