package middleware

import (
	"bytes"
	"io"
	"time"

	"emergency-his/server/logger"
	"github.com/gin-gonic/gin"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		start := time.Now()
		c.Next()
		logger.Business.Printf("request method=%s uri=%s query=%q body=%q status=%d duration=%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery, string(body), c.Writer.Status(), time.Since(start))
	}
}
