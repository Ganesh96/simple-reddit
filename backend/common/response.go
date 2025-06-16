package common

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status  int         `json:"status"`
	Message APIMessage  `json:"message"`
	Data    interface{} `json:"data"`
}

func NewAPIResponse(status int, message APIMessage, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func RespondWithJSON(c *gin.Context, status int, message APIMessage, data interface{}) {
	c.JSON(status, NewAPIResponse(status, message, data))
}
