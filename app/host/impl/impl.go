package impl

import (
	"database/sql"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/jacknotes/restful-api-demo/conf"
)

var Service *impl = &impl{}

type impl struct {
	// 可以更换成你们熟悉的，Logrus, 标准库log，zap
	// mcube Log模块是包装zap的实现
	log logger.Logger //记录日志
	// 依赖数据库
	db *sql.DB
}

func (i *impl) Init() error {
	i.log = zap.L().Named("Host")
	//获取全局db单例连接
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}
	i.db = db
	return nil
}
