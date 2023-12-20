package client

import (
	"context"

	"github.com/HumXC/shiroko/tools/input"
	"github.com/HumXC/shiroko/tools/minicap"
	"github.com/HumXC/shiroko/tools/screencap"
	"github.com/HumXC/shiroko/tools/shell"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Manager interface {
	// 返回所有工具的名字
	List() ([]string, error)
	// 对应 tools.common.Base

	Health(name string) error
	Install(name string) error
	Uninstall(name string) error
	Env(name string) ([]string, error)
	Exe(name string) (string, error)
	Args(name string) ([]string, error)
	Files(name string) ([]string, error)
}
type Minicap = minicap.IMinicap
type Screencap = screencap.IScreencap
type Input = input.IInput
type Shell = shell.IShell
type Client struct {
	conn      *grpc.ClientConn
	Screencap Screencap
	Minicap   Minicap
	Manager   Manager
	Input     Input
	Shell     Shell
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
		Input:     initInput(ctx, conn),
		Shell:     initShell(ctx, conn),
	}, nil
}
func (c *Client) Close() error {
	return c.conn.Close()
}
