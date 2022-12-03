package impl

import (
	"context"

	"github.com/jacknotes/restful-api-demo/app/host"
)

func (i *impl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	return nil, nil
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
