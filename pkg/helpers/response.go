// Package helpers
package helpers

import "github.com/gin-gonic/gin"

type Response struct {
	Error   bool   `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Errors  any    `json:"errors"`
}

func SuccessResponse(c *gin.Context, statusCode int, data any, message ...string) {
	response := Response{
		Error:   false,
		Status:  statusCode,
		Message: getMessageOrDefault(message, "OK"),
		Data:    data,
	}
	c.JSON(statusCode, response)
}

func ErrorResponse(c *gin.Context, statusCode int, message string, data ...any) {
	var errData any
	if len(data) > 0 {
		errData = data[0]
	} else {
		errData = map[string]string{}
	}

	response := Response{
		Error:   true,
		Status:  statusCode,
		Message: message,
		Errors:  errData,
	}
	c.JSON(statusCode, response)
}

func getMessageOrDefault(msgs []string, defaultMsg string) string {
	if len(msgs) > 0 && msgs[0] != "" {
		return msgs[0]
	}
	return defaultMsg
}
