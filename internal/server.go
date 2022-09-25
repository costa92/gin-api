package internal

import (
	"github.com/costa92/go-web/internal/db"
	genericapiserver "github.com/costa92/go-web/internal/server"
	"github.com/costa92/go-web/pkg/logger"
	"github.com/costa92/go-web/pkg/shutdown"
)

type apiServer struct {
	gs               *shutdown.GracefulShutdown
	genericAPIServer *genericapiserver.GenericAPIServer
	gRPCAPIServer    *grpcAPIServer
}

// 预运行
func (a *apiServer) preRun() preparedAPIServer {
	initRoute(a.genericAPIServer.Engine)

	a.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		if db.MysqlStorage != nil {
			db.Close()
		}
		a.genericAPIServer.Close()
		return nil
	}))
	return preparedAPIServer{a}
}

// 预运行构造体
type preparedAPIServer struct {
	*apiServer
}

// Run 运行服务
func (s preparedAPIServer) Run() error {
	if err := s.gs.Start(); err != nil {
		logger.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	return s.genericAPIServer.Run()
}
