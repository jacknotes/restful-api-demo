package impl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/jacknotes/restful-api-demo/app/host"
	"github.com/rs/xid"
)

func (i *impl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	// 校验数据的合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	// 生成全局id
	ins.Id = xid.New().String()
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
			if err != nil {
				i.log.Debugf("tx rollback error,%s", err)
			}
		} else {
			err := tx.Commit()
			if err != nil {
				i.log.Debugf("tx commit error, %s", err)
			}
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

	_, err = descStmt.Exec(ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec, ins.OSType, ins.OSName,
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

	// 用户输入了关键字
	// Prepre占位符必须是问号，"%kws%"是一个点位符
	if req.Keywords != "" {
		query.Where("r.name like ?", "%"+req.Keywords+"%")
	}

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
	// 新建一个sqlbuilder 查询对象，后面是语句
	query := sqlbuilder.NewQuery(queryHostSQL).Where("r.id = ?", req.Id)

	// 构建排序查询语句，并且会拼接成 id=?,name=?  的格式
	sqlStr, args := query.BuildQuery()
	i.log.Debugf("sql: %s, args: %s", sqlStr, args)

	// sql语句分析
	stmt, err := i.db.Prepare(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("prepare stmt query error, %s", err)
	}
	defer stmt.Close()

	ins := host.NewDefaultHost()
	// 查询时传入填充问号的值，上面BuildQuery已经加入了id=?,name=? 格式
	err = stmt.QueryRow(args...).Scan(
		&ins.Id, &ins.Vendor, &ins.Region, &ins.Zone, &ins.CreateAt, &ins.ExpireAt,
		&ins.Category, &ins.Type, &ins.InstanceId, &ins.Name,
		&ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.SyncAccount,
		&ins.PublicIP, &ins.PrivateIP, &ins.PayType, &ins.ResourceHash, &ins.DescribeHash,
		&ins.Id, &ins.CPU,
		&ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName,
		&ins.SerialNumber, &ins.ImageID, &ins.InternetMaxBandwidthOut, &ins.InternetMaxBandwidthIn,
		&ins.KeyPairName, &ins.SecurityGroups,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// 让返回code从200变成404
			return nil, exception.NewNotFound("host %s not found", req.Id)
		}
		return nil, fmt.Errorf("stmt query error, %s", err)
	}

	return ins, nil
}

// 需要实现事务，实现hosts和resource两个表的操作进行原子操作
func (i *impl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	var (
		resStmt  *sql.Stmt
		descStmt *sql.Stmt
		err      error
	)

	// 查询请求中的resource.Id是否存在
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithID(req.Id))
	if err != nil {
		return nil, err
	}

	// 对象更新（PATCH/PUT）
	switch req.UpdateMode {
	case host.PUT:
		// 对象更新（全量更新），存在于内存中，此时并未更新数据库
		ins.Update(req.Resource, req.Describe)
	case host.PATCH:
		// 对象打补丁（部分更新）
		err := ins.Patch(req.Resource, req.Describe)
		if err != nil {
			return nil, err
		}
	}

	// 校验更新后的数据是否合法，防止用户传入错误的数据
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	// 初始化一个事务,所有的操作都使用这个事务来进行提交，传入ctx上下方，如果上下文取消则整个事务操作都会取消
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 函数执行完成后专门判断事务是否正常
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				i.log.Debugf("tx rollback error,%s", err)
			}
		} else {
			err := tx.Commit()
			if err != nil {
				i.log.Debugf("tx commit error, %s", err)
			}
		}
	}()

	// 需要判断事务执行过程当中是否有异常
	// 事务有异常就回滚、无异常就提交
	// 在这个事务里面执行Insert SQL，为了防止sql注入攻击，需要使用PREPARE
	resStmt, err = tx.Prepare(updateResourceSQL)
	if err != nil {
		return nil, err
	}
	// 注意：Prepare语句会占用Mysql资源，一直不关闭会导致Prepare溢出，致使mysql全局报错无法使用，记得关掉
	defer resStmt.Close()

	// 将ins对象的属性更新到数据库中，这一步是真正改变数据库
	_, err = resStmt.Exec(ins.Vendor, ins.Region, ins.Zone, ins.ExpireAt, ins.CreateAt, ins.UpdateAt, ins.Name, ins.Description, ins.Id)
	if err != nil {
		return nil, err
	}

	descStmt, err = tx.Prepare(updateHostSQL)
	if err != nil {
		return nil, err
	}
	// 注意：Prepare语句会占用Mysql资源，一直不关闭会导致Prepare溢出，致使mysql全局报错无法使用，记得关掉
	defer descStmt.Close()

	// 将ins对象的属性更新到数据库中，这一步是真正改变数据库
	_, err = descStmt.Exec(ins.CPU, ins.Memory, ins.ResourceID)
	if err != nil {
		return nil, err
	}

	// debug
	i.log.Debugf("%s, %s, host.resource_id: %s, res.id: %s, %s", ins.CPU, ins.Memory, ins.ResourceID, ins.Id, ins.UpdateAt)
	return ins, nil
}

func (i *impl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {
	var (
		resStmt  *sql.Stmt
		descStmt *sql.Stmt
		err      error
	)

	// 查询请求中的resource.Id是否存在
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithID(req.Id))
	if err != nil {
		return nil, err
	}

	// 初始化一个事务,所有的操作都使用这个事务来进行提交，传入ctx上下方，如果上下文取消则整个事务操作都会取消
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 函数执行完成后专门判断事务是否正常
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				i.log.Debugf("tx rollback error,%s", err)
			}
		} else {
			err := tx.Commit()
			if err != nil {
				i.log.Debugf("tx commit error, %s", err)
			}
		}
	}()

	// 需要判断事务执行过程当中是否有异常
	// 事务有异常就回滚、无异常就提交
	// 在这个事务里面执行Insert SQL，为了防止sql注入攻击，需要使用PREPARE
	resStmt, err = tx.Prepare(deleteResourceSQL)
	if err != nil {
		return nil, err
	}
	// 注意：Prepare语句会占用Mysql资源，一直不关闭会导致Prepare溢出，致使mysql全局报错无法使用，记得关掉
	defer resStmt.Close()

	_, err = resStmt.Exec(req.Id)
	if err != nil {
		return nil, err
	}

	descStmt, err = tx.Prepare(deleteHostSQL)
	if err != nil {
		return nil, err
	}
	// 注意：Prepare语句会占用Mysql资源，一直不关闭会导致Prepare溢出，致使mysql全局报错无法使用，记得关掉
	defer descStmt.Close()

	_, err = descStmt.Exec(req.Id)
	if err != nil {
		return nil, err
	}
	return ins, nil
}
