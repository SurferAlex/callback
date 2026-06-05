package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLog logs method, path, status and latency for non-health routes.
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		path := c.Request.URL.Path
		if path == "/api/v1/ping" || path == "/api/v1/health" {
			return
		}
		if c.Writer.Status() >= 400 {
			log.Printf("[http] %s %s status=%d latency=%s client=%s",
				c.Request.Method, path, c.Writer.Status(), time.Since(start), c.ClientIP())
		}
	}
}
