package http

import (
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/jacknotes/restful-api-demo/apps/host"
	"github.com/julienschmidt/httprouter"
)

// Host模块的HTTP API 服务实例，暴露给最外层protocol层
var API = &handler{}

type handler struct {
	host host.ServiceServer // 最顶层抽象的接口，传入的值为实现此接口的对象
	log  logger.Logger
}

// 初始化的时候，依赖外部Host Service的实例对象
func (h *handler) Init() {
	// 创建日志子系统"HOST API"
	h.log = zap.L().Named("HOST API")

	// 获取grpc的对象，断言类型为host.ServiceServer
	h.host = app.GetGrpcApp(host.AppName).(host.ServiceServer)

}

func (h *handler) Name() string {
	return host.AppName
}

// 把handler实现的方法 注册给主路由
func (h *handler) Registry(r *httprouter.Router) {
	r.POST("/hosts", h.CreateHost)
	r.GET("/hosts", h.QueryHost)
	// restful API格式，路径匹配，路径参数，例如: /hosts/11001
	r.GET("/hosts/:id", h.DescribeHost)
	r.PUT("/hosts/:id", h.UpdateHost)
	r.PATCH("/hosts/:id", h.PatchHost)
	r.DELETE("/hosts/:id", h.DeleteHost)
}
