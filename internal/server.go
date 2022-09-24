package internal

import (
	genericapiserver "github.com/costa92/go-web/internal/server"
)

type apiServer struct {
	genericAPIServer *genericapiserver.GenericAPIServer
	gRPCAPIServer    *grpcAPIServer
}

// 预运行
func (a *apiServer) preRun() preparedAPIServer {
	initRoute(a.genericAPIServer.Engine)

	return preparedAPIServer{a}
}

// 预运行
type preparedAPIServer struct {
	*apiServer
}

func (s preparedAPIServer) Run() error {
	return s.genericAPIServer.Run()
}
