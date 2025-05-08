package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseIDParam(c *gin.Context, paramName string) (int, error) {
	idParam := c.Param(paramName)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Неверный формат ID", err)
		return 0, err
	}
	return id, nil
}
