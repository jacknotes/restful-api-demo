package impl

import (
	"context"
	"database/sql"

	"github.com/jacknotes/restful-api-demo/app/host"
)

func (i *impl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	// 校验数据的合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	// 把数据入库到resouce、host表
	// 一次需要往2个表录入数据，2个表必需同时存在数据或不存在数据，需要使用事务

	var (
		resStmt  *sql.Stmt
		descStmt *sql.Stmt
		err      error
	)

	// 初始化一个事务,所有的操作都使用这个事务来进行提交，传入ctx上下方，如果上下文取消则整个事务操作都会取消
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 函数执行完成后专门判断事务是否正常
	defer func() {
		if err != nil {
			err := tx.Rollback()
			i.log.Debugf("tx rollback error,%s", err)
		} else {
			err := tx.Commit()
			i.log.Debugf("tx commit error, %s", err)
		}
	}()

	// 需要判断事务执行过程当中是否有异常
	// 事务有异常就回滚、无异常就提交

	// 在这个事务里面执行Insert SQL，为了防止sql注入攻击，需要使用PREPARE
	resStmt, err = tx.Prepare(insertResourceSQL)
	if err != nil {
		return nil, err
	}
	// 注意：Prepare语句会占用Mysql资源，一直不关闭会导致Prepare溢出，致使mysql全局报错无法使用，记得关掉
	defer resStmt.Close()

	// 存储resource数据
	_, err = resStmt.Exec(ins.Id, ins.Vendor, ins.Region, ins.Zone, ins.CreateAt, ins.ExpireAt, ins.Category, ins.Type, ins.InstanceId,
		ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.SyncAccount, ins.PublicIP,
		ins.PrivateIP, ins.PayType, ins.ResourceHash, ins.DescribeHash,
	)
	if err != nil {
		return nil, err
	}

	// 存储host数据
	descStmt, err = tx.Prepare(insertDescribeSQL)
	if err != nil {
		return nil, err
	}
	defer descStmt.Close()

	_, err = resStmt.Exec(ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec, ins.OSType, ins.OSName,
		ins.SerialNumber, ins.ImageID, ins.InternetMaxBandwidthOut,
		ins.InternetMaxBandwidthIn, ins.KeyPairName, ins.SecurityGroups,
	)
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (i *impl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.Set, error) {
	return nil, nil
}

func (i *impl) DescribeHost(ctx context.Context, req *host.DescribeHostRequest) (*host.Host, error) {
	return nil, nil
}

func (i *impl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	return nil, nil
}

func (i *impl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {
	return nil, nil
}
