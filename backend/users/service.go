package users

import (
	"context"
	"net/http"
	"time"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // 14 is the cost factor
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password with a hashed password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser handles new user registration.
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user User

		if err := c.BindJSON(&user); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		// Check if username already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			common.RespondWithJSON(c, http.StatusConflict, common.USERNAME_ALREADY_EXISTS, gin.H{})
			return
		}

		// Hash the password before storing it
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, "Error hashing password", gin.H{"error": err.Error()})
			return
		}
		user.Password = hashedPassword

		// Insert the new user
		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		common.RespondWithJSON(c, http.StatusCreated, common.SUCCESS, gin.H{"user": result})
	}
}

// Login handles user authentication.
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user User
		var foundUser User

		if err := c.BindJSON(&user); err != nil {
			common.RespondWithJSON(c, http.StatusBadRequest, common.INVALID_REQUEST_BODY, gin.H{"error": err.Error()})
			return
		}

		// Find user by username
		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				common.RespondWithJSON(c, http.StatusUnauthorized, common.INVALID_CREDENTIALS, gin.H{})
				return
			}
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}

		// Check if the password is correct
		if !CheckPasswordHash(user.Password, foundUser.Password) {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.INVALID_CREDENTIALS, gin.H{})
			return
		}

		// Generate JWT token
		token, err := configs.GenerateJWT(foundUser.Username)
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, "Failed to generate token", gin.H{"error": err.Error()})
			return
		}

		common.RespondWithJSON(c, http.StatusOK, common.SUCCESS, gin.H{"token": token, "username": foundUser.Username})
	}
}

// DeleteUser handles the deletion of a user account.
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		username := c.Param("username")

		// Verify ownership before deletion.
		tokenUsername, exists := c.Get("username")
		if !exists || tokenUsername != username {
			common.RespondWithJSON(c, http.StatusForbidden, common.FORBIDDEN, gin.H{})
			return
		}

		result, err := userCollection.DeleteOne(ctx, bson.M{"username": username})
		if err != nil {
			common.RespondWithJSON(c, http.StatusInternalServerError, common.MONGO_DB_ERROR, gin.H{"error": err.Error()})
			return
		}
		if result.DeletedCount == 0 {
			common.RespondWithJSON(c, http.StatusNotFound, common.USER_NOT_FOUND, gin.H{})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
