package users

import (
	"net/http"
	"strings"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
)

// AuthorizeJWT is a middleware to authorize JWT tokens.
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const bearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.UNAUTHORIZED, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, bearerSchema) {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.UNAUTHORIZED, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := configs.ValidateToken(authHeader[len(bearerSchema):])
		if err != nil || !token.Valid {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.UNAUTHORIZED, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*configs.JWTClaim)
		if !ok || claims.Username == "" {
			common.RespondWithJSON(c, http.StatusUnauthorized, common.UNAUTHORIZED, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
