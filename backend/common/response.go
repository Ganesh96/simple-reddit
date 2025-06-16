package common

import "github.com/gin-gonic/gin"

// APIMessage is a generic struct for API responses
type APIMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// RespondWithJSON is a helper function to send a JSON response
func RespondWithJSON(c *gin.Context, httpStatus int, code string, payload interface{}) {
	var message string
	if msg, ok := SuccessMessages[code]; ok {
		message = msg.Message
	} else if err, ok := ErrorMessages[code]; ok {
		message = err.Message
	}

	c.JSON(httpStatus, gin.H{
		"status":  httpStatus,
		"message": message,
		"code":    code,
		"data":    payload,
	})
}
