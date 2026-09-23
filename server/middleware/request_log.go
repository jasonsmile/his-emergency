package middleware

import (
	"bytes"
	"encoding/json"
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
		logger.Business.Printf("request method=%s uri=%s query=%q body=%q status=%d duration=%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery, requestBodyForLog(body), c.Writer.Status(), time.Since(start))
		if c.Writer.Status() >= 400 {
			if _, logged := c.Get("error_logged"); !logged {
				logger.Error.Printf("request failed method=%s uri=%s status=%d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
			}
		}
	}
}

func requestBodyForLog(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err == nil {
		if _, exists := payload["password"]; exists {
			payload["password"] = "***"
		}
		if encoded, err := json.Marshal(payload); err == nil {
			return string(encoded)
		}
	}
	return string(body)
}
