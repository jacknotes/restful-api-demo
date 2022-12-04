package host

import "context"

type Service interface {
	// context.Context内部服务之间的调用链关系，
	// 录入主机信息
	CreateHost(context.Context, *Host) (*Host, error)
	// 查询主机列表信息
	QueryHost(context.Context, *QueryHostRequest) (*Set, error)
	// 查询主机详情
	DescribeHost(context.Context, *DescribeHostRequest) (*Host, error)
	// 修改主机信息
	UpdateHost(context.Context, *UpdateHostRequest) (*Host, error)
	// 删除主机， 为了兼容GRPC和 delete event需要返回Host
	DeleteHost(context.Context, *DeleteHostRequest) (*Host, error)
}

// 查询主机列表信息 传入参数
type QueryHostRequest struct {
	PageSize   int
	PageNumber int
}

func (req *QueryHostRequest) Offset() int {
	return (req.PageNumber - 1) * req.PageSize
}

// 查询主机详情 传入参数
type DescribeHostRequest struct {
	Id string
}

// 修改主机信息 传入参数
type UpdateHostRequest struct {
	Id string
}

// 删除主机 传入参数
type DeleteHostRequest struct {
	Id string
}
