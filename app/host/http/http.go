package http

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/jacknotes/restful-api-demo/app"
	"github.com/jacknotes/restful-api-demo/app/host"
	"github.com/julienschmidt/httprouter"
)

// Host模块的HTTP API 服务实例
var API = &handler{}

type handler struct {
	host host.Service
	log  logger.Logger
}

// 初始化的时候，依赖外部Host Service的实例对象
func (h *handler) Init() {
	h.log = zap.L().Named("HOST API")
	if app.Host == nil {
		panic("dependence host service is nil")
	}
	h.host = app.Host

}

// 把handler实现的方法 注册给主路由
func (h *handler) Registry(r *httprouter.Router) {
	r.POST("/hosts", h.CreateHost)
	r.GET("/hosts", h.QueryHost)
}