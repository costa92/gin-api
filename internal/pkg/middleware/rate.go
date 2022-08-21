package middleware

import (
	"strconv"
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

// ErrLimitExceeded defines Limit exceeded error.
var ErrLimitExceeded = errors.New("Limit exceeded")

func RateLimit() gin.HandlerFunc {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	limiter := redis_rate.NewLimiter(rdb)
	return func(ctx *gin.Context) {
		res, _ := limiter.Allow(ctx, "project:123", redis_rate.PerMinute(10))
		h := ctx.Writer.Header()
		h.Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))
		if res.Allowed == 0 {
			seconds := int(res.RetryAfter / time.Second)
			h.Set("RateLimit-RetryAfter", strconv.Itoa(seconds))

			// Stop processing and return the error.
			// Limit reached
			_ = ctx.Error(ErrLimitExceeded)
			ctx.AbortWithStatus(429)
		}
		ctx.Next()
	}
}
