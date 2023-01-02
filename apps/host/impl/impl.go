package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/jacknotes/restful-api-demo/apps/host"
	"github.com/jacknotes/restful-api-demo/conf"
	"google.golang.org/grpc"
)

// 暴露底层服务给HTTP层
var Service *impl = &impl{}

type impl struct {
	// 可以更换成你们熟悉的，Logrus, 标准库log，zap
	// mcube Log模块是包装zap的实现
	log logger.Logger //记录日志
	// 依赖数据库
	db *sql.DB
	// 结构体嵌套，继承grpc UnimplementedServiceServer对象，从而实现grpc ServiceServer接口
	host.UnimplementedServiceServer
}

func (i *impl) Config() error {
	i.log = zap.L().Named("Host")
	//获取全局db单例连接
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}
	i.db = db
	return nil
}

func (i *impl) Name() string {
	return host.AppName
}

func (i *impl) Registry(server *grpc.Server) {
	host.RegisterServiceServer(server, Service)
}

func init() {
	// grpc注册
	app.RegistryGrpcApp(Service)
}
