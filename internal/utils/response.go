package utils

import (
	"github.com/gin-gonic/gin"
)

type ResponseDTO struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}

type ErrorResponseDTO struct {
	Success bool `json:"success" example:"false"`
	Error   struct {
		Message string `json:"message" example:"Описание ошибки"`
		Debug   string `json:"debug,omitempty" example:"Детали ошибки для отладки"`
	} `json:"error"`
}

func Response(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := gin.H{
		"success": false,
		"error": gin.H{
			"message": message,
		},
	}

	if err != nil && gin.Mode() != gin.ReleaseMode {
		response["error"].(gin.H)["debug"] = err.Error()
	}

	c.JSON(statusCode, response)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidationErrorResponse(c *gin.Context, errors []ValidationError) {
	c.JSON(400, gin.H{
		"success": false,
		"error": gin.H{
			"message": "Ошибка валидации",
			"errors":  errors,
		},
	})
}
