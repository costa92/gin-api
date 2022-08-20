package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"time"
)

const TimeFieldFormat = "2006-01-02 15:04:00"

func Logger() gin.HandlerFunc {
	return LoggerWithWriter(gin.DefaultWriter)
}

func LoggerWithWriter(out io.Writer) gin.HandlerFunc {
	// 时间根式化
	zerolog.TimeFieldFormat = TimeFieldFormat
	zerolog.TimestampFunc = func() time.Time {
		// 时区
		return time.Now().In(time.Local)
	}
	log := zerolog.New(out).With().Timestamp().Logger()
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if raw != "" {
			path = path + "?" + raw
		}
		event := log.Info()
		if comment != "" {
			event = log.Error()
		}
		event.Int("statusCode", statusCode).
			Dur("latency", latency).
			Str("clientIP", clientIP).
			Str("method", method).
			Str("path", path).
			Msg(comment)
	}
}
