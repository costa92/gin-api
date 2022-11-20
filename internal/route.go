package internal

import (
	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/controller/v1"
	"github.com/costa92/go-web/controller/v1/areas"
	"github.com/costa92/go-web/controller/v1/distribute"
	enterprise2 "github.com/costa92/go-web/controller/v1/enterprise"
	"github.com/costa92/go-web/controller/v1/enterprise_contact"
	"github.com/costa92/go-web/controller/v1/enterprise_type"
	"github.com/costa92/go-web/controller/v1/follow"
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

	areaCtx := areas.NewAreaController(db.MysqlStorage)
	api.GET("/area/provinces", areaCtx.GetProvinces)
	api.GET("/area/areas", areaCtx.GetAreasByPid)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(c *gin.Context) {
		util.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
		return
	})

	api.Use(auto.AuthFunc())
	roleCtx := roles.NewRoleController(db.MysqlStorage)
	authCtx := v1.NewAuthController(db.MysqlStorage)
	menuCtx := menus.NewMenuController(db.MysqlStorage)
	userCtx := users.NewUserController(db.MysqlStorage)
	enterCtx := enterprise2.NewEnterpriseController(db.MysqlStorage)
	enterTypeCtx := enterprise_type.NewEnterpriseTypeController(db.MysqlStorage)
	enterContactCtx := enterprise_contact.NewEnterpriseContactController(db.MysqlStorage)
	distributeCtx := distribute.NewDistributeController(db.MysqlStorage)
	followCtx := follow.NewFollowController(db.MysqlStorage)

	api.GET("/getUserInfo", authCtx.GetUserInfo)
	api.GET("/getMenuList", menuCtx.GetMenuList)
	api.GET("/getPermCode", menuCtx.GetPermissionCode)
	api.GET("/enterprise/type/options", enterTypeCtx.GetOptions)

	api.GET("/users/options", userCtx.GetOptions)

	menu := api.Group("/")
	{
		menu.GET("menus", menuCtx.List)
		menu.POST("menu", menuCtx.Create)
		menu.PUT("menu", menuCtx.Update)
		menu.GET("menu", menuCtx.Detail)
	}
	role := api.Group("/")
	{
		role.GET("roles", roleCtx.Index)
		role.POST("role", roleCtx.Create)
		role.PUT("role", roleCtx.Update)
		role.GET("role", roleCtx.Detail)
		role.PUT("role/state", roleCtx.UpdateState)
		role.GET("role/all", roleCtx.GetAllList)
	}
	user := api.Group("/")
	{
		user.GET("users", userCtx.Users)
		user.POST("user", userCtx.Create)
		user.PUT("user", userCtx.Update)
		user.GET("user", userCtx.Get)
		user.PUT("user/state", userCtx.UpdateStates)
		user.POST("user/account/exist", userCtx.PostUserAccountExit)
	}
	enter := api.Group("/")
	{
		enter.GET("enterprises", enterCtx.GetList)
		enter.POST("enterprise", enterCtx.Create)
		enter.PUT("enterprise", enterCtx.Update)
		enter.GET("enterprise", enterCtx.Detail)
		enter.DELETE("enterprise", enterCtx.Deleted)
		enter.PUT("enterprise/status", enterCtx.UpdateStatus)

		enter.GET("enterprise/waiting/follows", enterCtx.GetWaitingFollowedList)
	}

	enterpriseR := api.Group("/enterprise/")
	{
		// 联系人
		enterpriseR.POST("contact", enterContactCtx.Create)
		enterpriseR.GET("contact/enterprises", enterContactCtx.GetListByEnterpriseId)
		enterpriseR.GET("contacts", enterContactCtx.GeList)
		enterpriseR.GET("contact", enterContactCtx.Detail)
		enterpriseR.PUT("contact", enterContactCtx.Update)

		// 类型
		enterpriseR.GET("types", enterTypeCtx.Index)
		enterpriseR.PUT("type", enterTypeCtx.Update)
		enterpriseR.POST("type", enterTypeCtx.Create)
		enterpriseR.GET("type", enterTypeCtx.Detail)
	}

	distributeR := api.Group("/")
	{
		distributeR.POST("distribute/area", distributeCtx.AreaCreate)
		distributeR.POST("distribute", distributeCtx.SingleCreate)
	}

	followR := api.Group("/")
	{
		followR.GET("follows", followCtx.GetFollowList)
	}

	return g
}
