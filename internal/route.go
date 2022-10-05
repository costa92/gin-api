package internal

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/controller/v1"
	"github.com/costa92/go-web/controller/v1/menus"
	"github.com/costa92/go-web/controller/v1/roles"
	"github.com/costa92/go-web/controller/v1/users"
	"github.com/costa92/go-web/internal/db"
	auth2 "github.com/costa92/go-web/internal/middleware/auth"
	"github.com/costa92/go-web/pkg/code"
	"github.com/costa92/go-web/pkg/util"
)

func initRoute(e *gin.Engine) {
	initMiddleware(e)
	initController(e)
}

func initMiddleware(e *gin.Engine) {
}

func initController(g *gin.Engine) *gin.Engine {
	api := g.Group("/crm_api")
	jwtStrategy, _ := newJWTAuth().(auth2.JWTStrategy)
	api.POST("/login", jwtStrategy.LoginHandler)
	api.POST("/logout", jwtStrategy.LogoutHandler)
	api.POST("/refresh", jwtStrategy.RefreshHandler)
	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		util.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	roleCtx := roles.NewRoleController(db.MysqlStorage)
	api.Use(auto.AuthFunc())

	authCtx := v1.NewAuthController(db.MysqlStorage)
	api.GET("/getUserInfo", authCtx.GetUserInfo)

	role := api.Group("/")
	{
		role.GET("roles", roleCtx.Index)
		role.POST("role", roleCtx.Create)
		role.PUT("role", roleCtx.Update)
		role.GET("role", roleCtx.Detail)
		role.PUT("role/state", roleCtx.UpdateState)
	}

	userCtx := users.NewUserController(db.MysqlStorage)
	user := api.Group("/")
	{
		user.GET("users", userCtx.Users)
		user.POST("user", userCtx.Create)
		user.PUT("user", userCtx.Update)
		user.GET("user", userCtx.Get)
		user.PUT("user/state", userCtx.UpdateStates)
	}

	menuCtx := menus.NewMenuController(db.MysqlStorage)
	menu := api.Group("/")
	{
		menu.GET("menus", menuCtx.List)
		menu.POST("menu", menuCtx.Create)
		menu.PUT("menu", menuCtx.Update)
	}
	return g
}
