package client

import (
	"github.com/jacknotes/restful-api-demo/apps/host"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	conf *Config
}

func NewClient(conf *Config) (*Client, error) {
	conn, err := grpc.Dial(
		conf.Addr,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(conf.Authentication),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
		conf: conf,
	}, nil
}

// Host Grpc服务的客户端
func (c *Client) Host() host.ServiceClient {
	return host.NewServiceClient(c.conn)
}
