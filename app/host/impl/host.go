package impl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/infraboard/mcube/sqlbuilder"
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
	// 新建一个sqlbuilder 查询对象，后面是语句
	query := sqlbuilder.NewQuery(queryHostSQL).Order("create_at").Desc().Limit(int64(req.Offset()), uint(req.PageSize))

	// 构建排序查询语句，并且会拼接成 id=?,name=?  的格式
	sqlStr, args := query.BuildQuery()
	i.log.Debugf("sql: %s, args: %s", sqlStr, args)
	// sql语句分析
	stmt, err := i.db.Prepare(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("prepare order stmt query error, %s", err)
	}
	defer stmt.Close()
	// 查询时传入填充问号的值，上面BuildQuery已经加入了id=?,name=? 格式
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("order stmt query error, %s", err)
	}

	// 初始化需要返回的对象
	set := host.NewSet()
	// 迭代查询表里的数据
	for rows.Next() {
		ins := host.NewDefaultHost()
		// 扫描行数据到对应变量中
		if err := rows.Scan(
			&ins.Id, &ins.Vendor, &ins.Region, &ins.Zone, &ins.CreateAt, &ins.ExpireAt,
			&ins.Category, &ins.Type, &ins.InstanceId, &ins.Name,
			&ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.SyncAccount,
			&ins.PublicIP, &ins.PrivateIP, &ins.PayType, &ins.ResourceHash, &ins.DescribeHash,
			&ins.Id, &ins.CPU,
			&ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName,
			&ins.SerialNumber, &ins.ImageID, &ins.InternetMaxBandwidthOut, &ins.InternetMaxBandwidthIn,
			&ins.KeyPairName, &ins.SecurityGroups,
		); err != nil {
			return nil, err
		}
		set.Add(ins)
	}

	// Count 查询总数量
	countStr, countArgs := query.BuildCount()
	countStmt, err := i.db.Prepare(countStr)
	if err != nil {
		return nil, fmt.Errorf("prepare count stmt query error, %s", err)
	}
	defer countStmt.Close()

	// 查询出来的值赋值给&set.Total
	if err := countStmt.QueryRow(countArgs...).Scan(&set.Total); err != nil {
		return nil, fmt.Errorf("count stmt query error, %s", err)
	}

	return set, nil
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
