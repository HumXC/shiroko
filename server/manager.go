package server

import (
	"context"

	"github.com/HumXC/shiroko/protos/common"
	pManager "github.com/HumXC/shiroko/protos/manager"
	"github.com/HumXC/shiroko/tools/manager"
)

type serverManager struct {
	// 嵌入此是为了确保所有定义的方法都被实现
	pManager.UnsafeManagerServer
	manager manager.IManager
}

// Args implements manager.ManagerServer.
func (s *serverManager) Args(ctx context.Context, req *pManager.NameRequest) (*pManager.ArgsResponse, error) {
	args, err := s.manager.Args(req.Name)
	if err != nil {
		return nil, err
	}
	return &pManager.ArgsResponse{Args: args}, nil
}

// Env implements manager.ManagerServer.
func (s *serverManager) Env(ctx context.Context, req *pManager.NameRequest) (*pManager.EnvResponse, error) {
	envs, err := s.manager.Env(req.Name)
	if err != nil {
		return nil, err
	}
	return &pManager.EnvResponse{Envs: envs}, nil
}

// Exe implements manager.ManagerServer.
func (s *serverManager) Exe(ctx context.Context, req *pManager.NameRequest) (*pManager.ExeResponse, error) {
	exe, err := s.manager.Exe(req.Name)
	if err != nil {
		return nil, err
	}
	return &pManager.ExeResponse{Exe: exe}, nil
}

// Files implements manager.ManagerServer.
func (s *serverManager) Files(ctx context.Context, req *pManager.NameRequest) (*pManager.FilesResponse, error) {
	files, err := s.manager.Files(req.Name)
	if err != nil {
		return nil, err
	}
	return &pManager.FilesResponse{Files: files}, nil
}

// Health implements manager.ManagerServer.
func (s *serverManager) Health(ctx context.Context, req *pManager.NameRequest) (*common.Empty, error) {
	err := s.manager.Health(req.Name)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

// Install implements manager.ManagerServer.
func (s *serverManager) Install(ctx context.Context, req *pManager.NameRequest) (*common.Empty, error) {
	err := s.manager.Install(req.Name)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

// List implements manager.ManagerServer.
func (s *serverManager) List(context.Context, *common.Empty) (*pManager.ListResponse, error) {
	list := s.manager.List()
	return &pManager.ListResponse{Names: list}, nil
}

// Uninstall implements manager.ManagerServer.
func (s *serverManager) Uninstall(ctx context.Context, req *pManager.NameRequest) (*common.Empty, error) {
	err := s.manager.Uninstall(req.Name)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

var _ pManager.ManagerServer = &serverManager{}

func NewManagerServer() *serverManager {
	return &serverManager{manager: manager.Manager}
}
