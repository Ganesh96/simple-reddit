package common

import "net/http"

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewAPIResponse(status int, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func RespondWithJSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, NewAPIResponse(status, message, data))
}