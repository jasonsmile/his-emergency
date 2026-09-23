package response

import (
	"net/http"

	"emergency-his/server/logger"
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

func SuccessWithPage(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": data, "total": total, "page": page, "page_size": pageSize})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: message, Data: data})
}

func Error(c *gin.Context, status, code int, message string) {
	logger.Error.Printf("request failed method=%s uri=%s status=%d code=%d message=%q", c.Request.Method, c.Request.URL.Path, status, code, message)
	c.Set("error_logged", true)
	c.JSON(status, Body{Code: code, Message: message})
}

func BadRequest(c *gin.Context, message string)   { Error(c, http.StatusBadRequest, 400, message) }
func Unauthorized(c *gin.Context, message string) { Error(c, http.StatusUnauthorized, 401, message) }
func NotFound(c *gin.Context, message string)     { Error(c, http.StatusNotFound, 404, message) }
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 500, message)
}
