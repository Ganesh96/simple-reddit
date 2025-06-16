package common

const (
	// Success messages
	SUCCESS APIMessage = "success"
)
// File: backend/common/response.go
package common

import "[github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)"

// APIResponse is the standard structure for all API responses.
type APIResponse struct {
	Status  int         `json:"status"`
	Message APIMessage  `json:"message"`
	Data    interface{} `json:"data"`
}

// NewAPIResponse creates a new APIResponse instance.
func NewAPIResponse(status int, message APIMessage, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// RespondWithJSON sends a structured JSON response.
func RespondWithJSON(c *gin.Context, status int, message APIMessage, data interface{}) {
	c.JSON(status, NewAPIResponse(status, message, data))
}