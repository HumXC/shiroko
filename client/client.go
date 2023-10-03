package client

import (
	"context"

	"github.com/HumXC/shiroko/tools/manager"
	"github.com/HumXC/shiroko/tools/minicap"
	"github.com/HumXC/shiroko/tools/screencap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn      *grpc.ClientConn
	Screencap screencap.IScreencap
	Minicap   minicap.IMinicap
	Manager   manager.IManager
}

func New(target string, opts ...grpc.DialOption) (*Client, error) {
	_opts := append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	conn, err := grpc.Dial(target, _opts...)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	return &Client{
		conn:      conn,
		Screencap: initScreencap(ctx, conn),
		Minicap:   initMinicap(ctx, conn),
		Manager:   initManager(ctx, conn),
	}, nil
}
func (c *Client) Close() error {
	return c.conn.Close()
}
