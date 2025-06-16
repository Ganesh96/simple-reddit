package users

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.UNAUTHORIZED, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BEARER_SCHEMA) {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.UNAUTHORIZED, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := configs.ValidateToken(tokenString)

		if err != nil {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.UNAUTHORIZED, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			c.Set("username", claims["username"])
			c.Next()
		} else {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.UNAUTHORIZED, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}
