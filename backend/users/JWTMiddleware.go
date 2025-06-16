package users

import (
	"net/http"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
)

// AuthorizeJWT is a middleware to protect routes that require authentication.
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.RespondWithJSON(c, http.StatusUnauthorized, "No auth token provided", nil)
			c.Abort()
			return
		}

		// The token is expected to be in the format "Bearer <token>"
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := configs.ValidateToken(tokenString)

		if token.Valid {
			// If the token is valid, extract the claims and set the username in the context
			// for subsequent handlers to use.
			claims := token.Claims.(*configs.JWTClaim)
			c.Set("username", claims.Username)
		} else {
			// If the token is invalid, respond with an unauthorized error.
			common.RespondWithJSON(c, http.StatusUnauthorized, "Invalid token", gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}
}
