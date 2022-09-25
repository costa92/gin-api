package server

import (
	"context"
	"net/http"
	"time"

	"github.com/costa92/errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/costa92/go-web/internal/metrics"
	"github.com/costa92/go-web/internal/middleware"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/util"
)

type GenericAPIServer struct {
	middlewares []string
	*gin.Engine
	SecureServingInfo            *SecureServingInfo
	InsecureServingInfo          *InsecureServingInfo
	healthz                      bool
	insecureServer, secureServer *http.Server
}

func initGenericAPIServer(s *GenericAPIServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

func (s *GenericAPIServer) InstallAPIs() {
	s.GET("/ping", func(c *gin.Context) {
		data := map[string]string{
			"message": "pong",
		}
		util.WriteResponse(c, nil, data)
	})
	if s.healthz {
		// 检查健康接口
		s.GET("/health", func(c *gin.Context) {
			data := map[string]string{
				"message": "health",
			}
			util.WriteResponse(c, nil, data)
		})
	}

	s.GET("/version", func(c *gin.Context) {
		data := map[string]string{
			"message": "v1",
		}
		util.WriteResponse(c, nil, data)
	})

	metrics.Metrics(s.Engine)
}

func (s *GenericAPIServer) Setup() {
	// 处理日志
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.Infow("DebugPrintRouteFunc ", "httpMethod", httpMethod, "absolutePath",
			absolutePath, "handlerName", handlerName, "nuHandlers", nuHandlers)
	}
}

// InstallMiddlewares 初始化中间件
func (s *GenericAPIServer) InstallMiddlewares() {
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())

	for _, m := range s.middlewares {
		mw, ok := middleware.Middleware[m]
		if !ok {
			continue
		}
		s.Use(mw)
	}
}

// Run 启动服务
func (s *GenericAPIServer) Run() error {
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServingInfo.Address, // 获取ip 与 端口号
		Handler: s,
	}

	s.secureServer = &http.Server{
		Addr:    s.SecureServingInfo.Address(),
		Handler: s,
	}

	var eg errgroup.Group
	// 运行 http
	eg.Go(func() error {
		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalw("启动失败", "err", err)
			return err
		}
		return nil
	})

	// 运行 https
	eg.Go(func() error {
		key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
		if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
			return nil
		}
		logger.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())
		if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalw(err.Error())
			return err
		}
		logger.Infof("Server on %s stopped", s.SecureServingInfo.Address())
		return nil
	})

	if err := eg.Wait(); err != nil {
		logger.Fatalw("err group Wait", "err", err)
	}
	return nil
}

func (s *GenericAPIServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.insecureServer.Shutdown(ctx); err != nil {
		logger.Warnf("Shutdown insecure server failed: %s", err.Error())
	}
}
