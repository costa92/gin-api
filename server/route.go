package server

import (
	"github.com/costa92/go-web/controller/auth"
	"github.com/gin-gonic/gin"

	"github.com/costa92/go-web/controller/index"
	"github.com/costa92/go-web/internal/db"
)

func initRoute(e *gin.Engine) {
	initMiddleware(e)
	initController(e)
}

func initMiddleware(e *gin.Engine) {
}

func initController(g *gin.Engine) *gin.Engine {

	api := g.Group("/api")

	authCtx := auth.NewAuth(db.MysqlStorage)
	api.POST("login", authCtx.Login)

	index := index.NewIndex(db.MysqlStorage)
	// idx := g.Group("/index", middleware.RateLimit())
	idx := api.Group("/")
	{
		idx.GET("index", index.Index)
		idx.POST("index", index.Create)
		idx.PUT("index", index.Update)
		idx.GET("index/detail", index.Detail)
	}
	return g
}
