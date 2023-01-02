package protocol

import (
	"net"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	app "github.com/jacknotes/restful-api-demo/apps"
	"github.com/jacknotes/restful-api-demo/apps/host"
	"github.com/jacknotes/restful-api-demo/conf"
	"google.golang.org/grpc"
)

// GrpcService grpc服务
type GrpcService struct {
	// server 对象
	server *grpc.Server
	// 日志
	l logger.Logger
}

func NewGrpcService() *GrpcService {
	return &GrpcService{
		server: grpc.NewServer(),
		l:      zap.L().Named("GRPC Server"),
	}
}

func (s *GrpcService) Start() {
	// 加载服务
	host.RegisterServiceServer(s.server, app.Host)

	lis, err := net.Listen("tcp", conf.C().App.GrpcAddr())
	if err != nil {
		s.l.Errorf("listen grpc tcp conn error, %s", err)
		return
	}

	s.l.Infof("GRPC 服务监听地址: %s", conf.C().App.GrpcAddr())
	if err := s.server.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}

		s.l.Error("start grpc service error, %s", err.Error())
		return
	}
}

func (s *GrpcService) Stop() {
	s.server.GracefulStop()
}
