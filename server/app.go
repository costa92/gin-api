package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/costa92/go-web/config"
	"github.com/costa92/go-web/internal/metrics"
	"github.com/costa92/go-web/internal/pkg/middleware"
	"github.com/costa92/go-web/internal/pkg/util"
)

type Server struct {
	Conf        *config.ServerConf
	middlewares []string
}

func NewServer(conf *config.ServerConf) *Server {
	return &Server{
		Conf:        conf,
		middlewares: conf.Middlewares,
	}
}

// 预运行
func (sr *Server) preRun() *gin.Engine {
	gin.DisableConsoleColor()
	e := gin.Default()
	// 初始化中间件
	sr.InstallMiddlewares(e)
	sr.InstallAPIs(e)
	initRoute(e)
	return e
}

func (sr *Server) InstallAPIs(e *gin.Engine) {
	e.GET("/ping", func(c *gin.Context) {
		data := map[string]string{
			"message": "pong",
		}
		util.WriteResponse(c, nil, data)
	})
	// 检查健康接口
	e.GET("/health", func(c *gin.Context) {
		data := map[string]string{
			"message": "health",
		}
		util.WriteResponse(c, nil, data)
	})

	e.GET("/version", func(c *gin.Context) {
		data := map[string]string{
			"message": "v1",
		}
		util.WriteResponse(c, nil, data)
	})

	metrics.Metrics(e)
}

// Run  运行接口
func (sr *Server) Run() error {
	e := sr.preRun()
	srv := &http.Server{
		Addr:    ":" + sr.Conf.Port,
		Handler: e,
	}
	var eg errgroup.Group

	eg.Go(func() error {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
			return err
		}
		log.Fatalf("Server on %s stopped", srv.Addr)
		return nil
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}
	// 等待中断信号以优雅地关闭服务器(设置 5 秒钟的超时)
	quit := make(chan os.Signal)

	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
	return nil
}

// InstallMiddlewares 初始化中间件
func (sr *Server) InstallMiddlewares(e *gin.Engine) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Context())
	for _, m := range sr.middlewares {
		mw, ok := middleware.Middleware[m]
		if !ok {
			continue
		}
		e.Use(mw)
	}
}
