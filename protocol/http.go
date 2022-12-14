package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/jacknotes/restful-api-demo/conf"
)

// HTTPService http服务
type HTTPService struct {
	// router, root router，路由， method+path --> handler
	r router.Router
	// 日志
	l logger.Logger
	// c      *conf.Config
	// 服务实例对象，HTTP服务器
	server *http.Server
}

func NewHTTPService() *HTTPService {
	r := httprouter.New()
	// 开启API URI路径查看
	r.EnableAPIRoot()

	return &HTTPService{
		r: r,
		l: zap.L().Named("HTTP Server"),
		server: &http.Server{
			// http server 监听地址
			Addr: conf.C().App.Addr(),
			// http handler/router
			// httprouter对象实现了Handler接口
			Handler: r,
			// 服务端读取Header 超时设置
			ReadHeaderTimeout: 60 * time.Second,
			// 连接，client --> server，服务端读取超时时间
			ReadTimeout: 60 * time.Second,
			// 服务端响应超时时间
			WriteTimeout: 60 * time.Second,
			// 空闲超时时间，长连接
			IdleTimeout: 60 * time.Second,
			// 最大Header大小
			MaxHeaderBytes: 1 << 20, // 1M
		},
	}
}

// Start 启动服务
func (s *HTTPService) Start() error {
	// 装置子服务路由

	// hostAPI.API.Init() //此时初始化读取的h.host = app.Host = impl.Service
	// hostAPI.API.Registry(s.r)

	// 完成http app的注册
	app.LoadHttpApp("", s.r)

	// 启动 HTTP服务
	s.l.Infof("HTTP服务启动成功, 监听地址: %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}

// Stop 停止server
func (s *HTTPService) Stop() error {
	s.l.Info("start graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// 优雅关闭HTTP服务
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Errorf("graceful shutdown timeout, force exit")
	}
	return nil
}
