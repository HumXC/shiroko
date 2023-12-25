package server

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/HumXC/shiroko/protos/common"
	pShell "github.com/HumXC/shiroko/protos/shell"
	"github.com/HumXC/shiroko/tools/shell"
)

// TODO 优化这玩意
type ReadWriteCloser struct {
	buf    bytes.Buffer
	isDone chan struct{}
	mu     sync.Mutex
}

func NewRWCloser() *ReadWriteCloser {
	return &ReadWriteCloser{
		buf:    bytes.Buffer{},
		isDone: make(chan struct{}),
		mu:     sync.Mutex{},
	}
}
func (rwc *ReadWriteCloser) Close() error {
	select {
	case <-rwc.isDone:
	default:
		close(rwc.isDone)
	}
	return nil
}
func (rwc *ReadWriteCloser) Read(p []byte) (int, error) {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	select {
	case <-rwc.isDone:
		if rwc.buf.Len() == 0 {
			return 0, io.EOF
		}
		n, _ := rwc.buf.Read(p)
		return n, nil
	default:
		n, _ := rwc.buf.Read(p)
		return n, nil
	}
}
func (rwc *ReadWriteCloser) Write(p []byte) (int, error) {
	rwc.mu.Lock()
	defer rwc.mu.Unlock()
	select {
	case <-rwc.isDone:
		return 0, io.ErrClosedPipe
	default:
		n, _ := rwc.buf.Write(p)
		return n, nil
	}
}

type serverShell struct {
	// 嵌入此是为了确保所有定义的方法都被实现
	pShell.UnsafeShellServer
	shell shell.IShell
}

// Pull implements shell.ShellServer.
func (s *serverShell) Pull(req *pShell.PullRequest, serv pShell.Shell_PullServer) error {
	data, err := s.shell.Pull(req.Filename)
	if err != nil {
		return err
	}
	defer data.Close()
	buf := make([]byte, 1024)
	for {
		n, err := data.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if err != nil {
			return err
		}
		err = serv.Send(&common.DataChunk{Data: buf[:n]})
		if err != nil {
			return err
		}
	}
	return nil
}

// Push implements shell.ShellServer.
func (s *serverShell) Push(serv pShell.Shell_PushServer) error {
	var err error
	ready := make(chan struct{}, 1)
	defer close(ready)
	filename := ""
	buf := NewRWCloser()
	defer buf.Close()
	go func(buf io.WriteCloser, ser pShell.Shell_PushServer) {
		isSetFilename := false
		for {
			req, _err := serv.Recv()
			if _err == io.EOF {
				buf.Close()
				break
			}
			if _err != nil {
				err = _err
				return
			}
			// ready 用于表示已经读取了 filename
			// 如果已经读取了 filename, 则写入 ready
			if !isSetFilename {
				isSetFilename = true
				filename = req.Filename
				ready <- struct{}{}
			}
			buf.Write(req.Data)
		}
	}(buf, serv)

	select {
	case <-ready:
		break
	case <-serv.Context().Done():
		return serv.Context().Err()
	}
	err = s.shell.Push(filename, buf)
	if err != nil {
		return err
	}
	err = serv.SendAndClose(&common.Empty{})
	return err
}

// Getprop implements shell.ShellServer.
func (s *serverShell) Getprop(ctx context.Context, req *pShell.GetpropRequest) (*pShell.GetpropResponse, error) {
	val, err := s.shell.Getprop(req.Key)
	return &pShell.GetpropResponse{Value: val}, err
}

// Install implements shell.ShellServer.
func (s *serverShell) Install(ctx context.Context, req *pShell.InstallRequest) (*common.Empty, error) {
	err := s.shell.Install(req.Apkpath)
	return &common.Empty{}, err
}

// ListApps implements shell.ShellServer.
func (s *serverShell) ListApps(context.Context, *common.Empty) (*pShell.ListAppsResponse, error) {
	result, err := s.shell.ListApps()
	return &pShell.ListAppsResponse{Apps: result}, err
}

// Run implements shell.ShellServer.
func (s *serverShell) Run(ctx context.Context, req *pShell.RunRequest) (*common.DataChunk, error) {
	out, err := s.shell.Run(req.Cmd, req.TimeoutMs)
	return &common.DataChunk{Data: out}, err
}

// StartApp implements shell.ShellServer.
func (s *serverShell) StartApp(ctx context.Context, req *pShell.StartAppRequest) (*common.Empty, error) {
	err := s.shell.StartApp(req.Activity)
	return &common.Empty{}, err
}

// StopApp implements shell.ShellServer.
func (s *serverShell) StopApp(ctx context.Context, req *pShell.StopAppRequest) (*common.Empty, error) {
	err := s.shell.StopApp(req.Pkgname)
	return &common.Empty{}, err
}

// Uninstall implements shell.ShellServer.
func (s *serverShell) Uninstall(ctx context.Context, req *pShell.UninstallRequest) (*common.Empty, error) {
	err := s.shell.Uninstall(req.Pkgname)
	return &common.Empty{}, err
}

var _ pShell.ShellServer = &serverShell{}

func NewShellServer() *serverShell {
	return &serverShell{
		shell: shell.Shell,
	}
}
