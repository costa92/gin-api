package middleware

import (
	"github.com/costa92/go-web/internal/logger"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

const TimeFieldFormat = "2006-01-02 15:04:00"

func Logger() gin.HandlerFunc {
	return LoggerWithWriter(gin.DefaultWriter)
}

func LoggerWithWriter(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if raw != "" {
			path = path + "?" + raw
		}
		logger.Infow("LoggerWithWriter",
			"statusCode", c.Writer.Status(),
			"latency", time.Since(start),
			"clientIP", c.ClientIP(),
			"method", c.Request.Method,
			"path", path,
			"comment", comment)
	}
}
