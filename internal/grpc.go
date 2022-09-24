package internal

import (
	"net"

	"google.golang.org/grpc"

	"github.com/costa92/go-web/pkg/logger"
)

type grpcAPIServer struct {
	*grpc.Server
	address string
}

func (s *grpcAPIServer) Run() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		logger.Fatalf("failed to listen: %s", err.Error())
	}

	go func() {
		if err := s.Serve(listen); err != nil {
			logger.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	logger.Infof("start grpc server at %s", s.address)
}

func (s *grpcAPIServer) Close() {
	s.GracefulStop()
	logger.Infof("GRPC server on %s stopped", s.address)
}
