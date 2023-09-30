package client

import (
	"context"

	"github.com/HumXC/shiroko/tools/screencap"
	"google.golang.org/grpc"
)

type Client struct {
	conn      *grpc.ClientConn
	Screencap screencap.IScreencap
}

func New(target string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	return &Client{
		conn:      conn,
		Screencap: initScreencap(ctx, conn),
	}, nil
}
func (c *Client) Close() error {
	return c.conn.Close()
}
