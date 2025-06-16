package users

import (
	"net/http"

	"github.com/ganesh96/simple-reddit/backend/common"
	"github.com/ganesh96/simple-reddit/backend/configs"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.RespondWithJSON(c, http.StatusUnauthorized, "No auth token provided", nil)
			c.Abort()
			return
		}

		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := configs.ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(*configs.JWTClaim)
			c.Set("username", claims.Username)
		} else {
			common.RespondWithJSON(c, http.StatusUnauthorized, "Invalid token", gin.H{"error": err.Error()})
			c.Abort()
			return
		}
	}
}
