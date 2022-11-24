package middleware

import (
	"github.com/gin-gonic/gin"
)

type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

type AuthOperator struct {
	strategy AuthStrategy
}

func (operator *AuthOperator) SetStrategy(strategy AuthStrategy) {
	operator.strategy = strategy
}

func (operator *AuthOperator) AuthFunc() gin.HandlerFunc {
	return operator.strategy.AuthFunc()
}

type AuthUser struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
}

func GetAuthUser(ctx *gin.Context) *AuthUser {
	authUser, _ := ctx.Get(UsernameKey)
	return authUser.(*AuthUser)
}

func GetAuthUserName(ctx *gin.Context) string {
	authUser := GetAuthUser(ctx)
	return authUser.Username
}

func GetAuthUserId(ctx *gin.Context) int {
	authUser := GetAuthUser(ctx)
	return authUser.UserId
}
