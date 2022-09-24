package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/pkg/logger"
)

const (
	UsernameKey = "username"
	UserIdKey   = "user_id"
)

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(logger.KeyUsername, c.GetString(UsernameKey))
		c.Next()
	}
}
