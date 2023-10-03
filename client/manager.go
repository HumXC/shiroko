package client

import (
	"context"

	pCommon "github.com/HumXC/shiroko/protos/common"
	pManager "github.com/HumXC/shiroko/protos/manager"
	tManager "github.com/HumXC/shiroko/tools/manager"
	"google.golang.org/grpc"
)

type managerClient struct {
	mm  pManager.ManagerClient
	ctx context.Context
}

// Args implements manager.IManager.
func (m *managerClient) Args(name string) ([]string, error) {
	resp, err := m.mm.Args(m.ctx, &pManager.NameRequest{Name: name})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Args, nil
}

// Env implements manager.IManager.
func (m *managerClient) Env(name string) ([]string, error) {
	resp, err := m.mm.Env(m.ctx, &pManager.NameRequest{Name: name})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Envs, nil
}

// Exe implements manager.IManager.
func (m *managerClient) Exe(name string) (string, error) {
	resp, err := m.mm.Exe(m.ctx, &pManager.NameRequest{Name: name})
	if err != nil {
		return "", ParseError(err)
	}
	return resp.Exe, nil
}

// Files implements manager.IManager.
func (m *managerClient) Files(name string) ([]string, error) {
	resp, err := m.mm.Files(m.ctx, &pManager.NameRequest{Name: name})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Files, nil
}

// Health implements manager.IManager.
func (m *managerClient) Health(name string) error {
	_, err := m.mm.Health(m.ctx, &pManager.NameRequest{Name: name})
	return ParseError(err)
}

// Install implements manager.IManager.
func (m *managerClient) Install(name string) error {
	_, err := m.mm.Install(m.ctx, &pManager.NameRequest{Name: name})
	return ParseError(err)
}

// List implements manager.IManager.
func (m *managerClient) List() []string {
	resp, err := m.mm.List(m.ctx, &pCommon.Empty{})
	err = ParseError(err) // 此处只能是 grpc 错误，考虑使用日志警告
	return resp.Names
}

// Uninstall implements manager.IManager.
func (m *managerClient) Uninstall(name string) error {
	_, err := m.mm.Uninstall(m.ctx, &pManager.NameRequest{Name: name})
	return ParseError(err)
}

func initManager(ctx context.Context, conn *grpc.ClientConn) tManager.IManager {
	return &managerClient{
		mm:  pManager.NewManagerClient(conn),
		ctx: ctx,
	}
}
