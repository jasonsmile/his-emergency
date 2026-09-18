package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "success", Data: data})
}

func Error(c *gin.Context, status, code int, message string) {
	c.JSON(status, Body{Code: code, Message: message})
}

func BadRequest(c *gin.Context, message string) { Error(c, http.StatusBadRequest, 400, message) }
func NotFound(c *gin.Context, message string)   { Error(c, http.StatusNotFound, 404, message) }
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 500, message)
}
