package users

import (
	"context"
	"net/http"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var user common.User
	if err := c.ShouldBindJSON(&user); err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
		return
	}

	// Check for existing user by email
	count, err := userCollection.CountDocuments(context.TODO(), bson.M{"email": user.Email})
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Error checking for existing email"})
		return
	}
	if count > 0 {
		common.RespondWithJSON(c, http.StatusConflict, common.EMAIL_ALREADY_EXISTS, gin.H{"error": "User with this email already exists"})
		return
	}

	// Check for existing user by username
	count, err = userCollection.CountDocuments(context.TODO(), bson.M{"username": user.Username})
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Error checking for existing username"})
		return
	}
	if count > 0 {
		common.RespondWithJSON(c, http.StatusConflict, common.USERNAME_ALREADY_EXISTS, gin.H{"error": "User with this username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()

	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to create user"})
		return
	}

	common.RespondWithJSON(c, http.StatusCreated, common.CREATED, gin.H{"message": "User created successfully"})
}

func Login(c *gin.Context) {
	var loginDetails LoginDetails
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
		return
	}

	var foundUser common.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": loginDetails.Email}).Decode(&foundUser)
	if err != nil {
		common.RespondWithJSON(c, http.StatusUnauthorized, common.INVALID_CREDENTIALS, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginDetails.Password))
	if err != nil {
		common.RespondWithJSON(c, http.StatusUnauthorized, common.INVALID_CREDENTIALS, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := configs.GenerateToken(foundUser.Username)
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to generate token"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"token": token})
}

func DeleteUser(c *gin.Context) {
	username := c.Param("username")
	authUsername := c.GetString("username")

	if username != authUsername {
		common.RespondWithJSON(c, http.StatusForbidden, common.FORBIDDEN, gin.H{"error": "You are not authorized to delete this user"})
		return
	}

	_, err := userCollection.DeleteOne(context.TODO(), bson.M{"username": username})
	if err != nil {
		common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": "Failed to delete user"})
		return
	}

	common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"message": "User deleted successfully"})
}
