package auth

import (
	"fmt"
	"time"

	"github.com/costa92/errors"
	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

var (
	ErrMissingKID    = errors.New("Invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("Can not obtain secret information from cache")
)

type Secret struct {
	Username string
	ID       string
	Key      string
	Expires  int64
}

type CacheStrategy struct {
	get func(kid string) (Secret, error)
}

var _ middleware.AuthStrategy = &CacheStrategy{}

func NewCacheStrategy(get func(kid string) (Secret, error)) CacheStrategy {
	return CacheStrategy{get}
}

func (cache CacheStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if len(header) == 0 {
			util.WriteResponse(
				c,
				errors.WithCode(code.ErrMissingHeader, "Authorization header cannot be empty."),
				nil,
			)
			c.Abort()
			return
		}
		var rawJWT string
		var secret Secret
		claims := &jwt.MapClaims{}
		parsedT, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is HMAC signature
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, ErrMissingKID
			}
			var err error
			secret, err = cache.get(kid)
			if err != nil {
				return nil, ErrMissingSecret
			}
			return []byte(secret.Key), nil
		}, jwt.WithAudience(AuthzAudience))

		if err != nil || !parsedT.Valid {
			util.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, err.Error()), nil)
			c.Abort()
			return
		}
		if KeyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			util.WriteResponse(c, errors.WithCode(code.ErrExpired, "expired at: %s", tm), nil)
			c.Abort()

			return
		}
		c.Set(middleware.UsernameKey, secret.Username)
		c.Next()
	}
}

func KeyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0))
	}
	return false
}
