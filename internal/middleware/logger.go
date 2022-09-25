package middleware

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"

	"github.com/costa92/go-web/pkg/logger"
)

const TimeFieldFormat = "2006-01-02 15:04:00"

var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
		// param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.ClientIP,
		param.Latency,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(GetLoggerConfig(nil, nil, nil))
}

func LoggerWithFormatter(f gin.LogFormatter) gin.HandlerFunc {
	return LoggerWithConfig(gin.LoggerConfig{
		Formatter: f,
	})
}

func LoggerWithWriter(out io.Writer, notLogger ...string) gin.HandlerFunc {
	return LoggerWithConfig(gin.LoggerConfig{
		Output:    out,
		SkipPaths: notLogger,
	})
}
func LoggerWithConfig(conf gin.LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}
	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notLogged := conf.SkipPaths
	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	if isTerm {
		gin.ForceConsoleColor()
	}
	var skip map[string]struct{}
	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}
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
