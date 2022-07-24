package middleware

import "github.com/gin-gonic/gin"

const UsernameKey = "username"

// Context is a middleware that injects common prefix fields to gin.Context.
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
