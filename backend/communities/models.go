package communities

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

var CommunitiesCollection *mongo.Collection = configs.GetCollection(configs.DB, "communities")

// Community struct
type Community struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name,omitempty"`
	Description   string             `bson:"description,omitempty"`
	Creation_date time.Time          `bson:"creation_date,omitempty"`
	Updation_date time.Time          `bson:"updation_date,omitempty"`
	Members_count int                `bson:"members_count,omitempty"`
	Posts_count   int                `bson:"posts_count,omitempty"`
	Creator_name  string             `bson:"creator_name,omitempty"`
}

type Communities []Community

// CreateCommunity creates a new community
func CreateCommunity(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var community Community

	if err := c.BindJSON(&community); err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
		return
	}

	newCommunity := Community{
		Id:            primitive.NewObjectID(),
		Name:          community.Name,
		Description:   community.Description,
		Creator_name:  community.Creator_name,
		Creation_date: time.Now(),
		Updation_date: time.Now(),
		Members_count: 1,
		Posts_count:   0,
	}

	result, err := CommunitiesCollection.InsertOne(ctx, newCommunity)
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
		return
	}
	common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"community": result})
}

// GetAllCommunities gets all communities
func GetAllCommunities(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var communities []Community

	cursor, err := CommunitiesCollection.Find(ctx, bson.M{})
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

// GetCommunityByName gets a single community by name
func GetCommunityByName(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var community Community
	communityName := c.Param("name")

	err := CommunitiesCollection.FindOne(ctx, bson.M{"name": communityName}).Decode(&community)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			common.RespondWithJSON(c, http.StatusNotFound, common.COMMUNITY_NOT_FOUND, nil)
			return
		}
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
		return
	}
	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"community": community})
}
