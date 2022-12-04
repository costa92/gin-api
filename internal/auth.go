package internal

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/costa92/go-web/internal/db"
	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/internal/middleware/auth"
	"github.com/costa92/go-web/model"
	"github.com/costa92/go-web/pkg/logger"
)

const (
	// APIServerAudience defines the value of jwt audience field.
	APIServerAudience = "go-web"

	// APIServerIssuer defines the value of jwt issuer field.
	APIServerIssuer = "apiserver"
)

type LoginInfo struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username, password string) bool {
		// fetch user from database
		userModel := model.NewUserModel(context.TODO(), db.MysqlStorage)
		user, err := userModel.FirstByName(username)
		if err != nil {
			return false
		}
		// Compare the login password with the user password.
		if err := user.Compare(password); err != nil {
			return false
		}
		return true
	})
}

func newJWTAuth() middleware.AuthStrategy {
	ginjwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"),
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          viper.GetDuration("jwt.timeout"),
		MaxRefresh:       viper.GetDuration("jwt.max-refresh"),
		Authenticator:    authenticator(),
		LoginResponse:    loginResponse(),
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": "Logout Success",
				"result":  map[string]string{},
			})
			return
		},
		RefreshResponse: refreshResponse(),
		PayloadFunc:     payloadFunc(),
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": message,
				"result":  map[string]string{},
			})
			return
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			userIdFloat64 := claims[middleware.UserIdKey].(float64)
			return &middleware.AuthUser{
				Username: claims[jwt.IdentityKey].(string),
				UserId:   int(userIdFloat64),
			}
		},
		IdentityKey:   middleware.UsernameKey,
		Authorizator:  authorizator(),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,
		TimeFunc:      time.Now,
	})
	return auth.NewJWTStrategy(*ginjwt)
}

func newAutoAuth() middleware.AuthStrategy {
	return auth.NewAutoStrategy(newBasicAuth().(auth.BasicStrategy), newJWTAuth().(auth.JWTStrategy))
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login LoginInfo
		var err error
		if c.Request.Header.Get("Authorization") != "" {
			login, err = parseWithHeader(c)
		} else {
			login, err = parseWithBody(c)
		}
		if err != nil {
			logger.Errorw("authenticator filed", "err", err)
			return "", jwt.ErrFailedAuthentication
		}
		userModel := model.NewUserModel(c, db.MysqlStorage)
		user, err := userModel.FirstByName(login.Username)
		if err != nil {
			logger.Errorw("get user information failed", "err", err)
			return "", jwt.ErrFailedAuthentication
		}
		if user.Status == model.StatusDisable {
			logger.Errorw("get user statue failed", "err", err)
			return "", errors.New("account status disable")
		}
		if err := user.Compare(login.Password); err != nil {
			logger.Errorf("user.Compare: %v", err)
			return "", jwt.ErrFailedAuthentication
		}
		user.LastTime = time.Now().Unix()
		if err := userModel.Save(user); err != nil {
			logger.Errorw("get user information failed", "err", err)
			return "", jwt.ErrFailedAuthentication
		}
		c.Set(middleware.UsernameKey, user.Username)
		c.Set(middleware.UserIdKey, user.ID)
		return user, nil
	}
}

func parseWithHeader(c *gin.Context) (LoginInfo, error) {
	author := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(author) != 2 || author[0] != "Basic" {
		logger.Errorw("get basic string from Authorization header failed")
		return LoginInfo{}, jwt.ErrFailedAuthentication
	}
	payload, err := base64.StdEncoding.DecodeString(author[1])
	if err != nil {
		logger.Errorw("decode basic string", "err", err)
		return LoginInfo{}, jwt.ErrFailedAuthentication
	}
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		logger.Errorw("parse payload failed")
		return LoginInfo{}, jwt.ErrFailedAuthentication
	}
	return LoginInfo{
		Username: pair[0],
		Password: pair[1],
	}, nil
}

func parseWithBody(c *gin.Context) (LoginInfo, error) {
	login := LoginInfo{}
	logger.Infow("parseWithBody start")
	if err := c.ShouldBindJSON(&login); err != nil {
		logger.Errorw("parseWithBody parse login parameters", "err", err)
		return LoginInfo{}, jwt.ErrFailedAuthentication
	}
	return login, nil
}

func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": "success",
			"result": map[string]string{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			},
		})
	}
}

func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		userId, _ := c.Get(middleware.UserIdKey)
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": "success",
			"result": map[string]interface{}{
				"token":    token,
				"expire":   expire.Format(time.RFC3339),
				"nickname": c.GetString(middleware.UsernameKey),
				"user_id":  userId,
			},
		})
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*middleware.AuthUser); ok {
			logger.Infof("user_id: `%d`, username: `%s` is authenticated.", v.UserId, v.Username)
			return true
		}
		return false
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		claims := jwt.MapClaims{
			"iss": APIServerIssuer,
			"aud": APIServerAudience,
		}
		if u, ok := data.(*model.User); ok {
			claims[jwt.IdentityKey] = u.Username
			claims[middleware.UserIdKey] = u.ID
			claims["sub"] = u.Username
		}
		return claims
	}
}
