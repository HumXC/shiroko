package client

import (
	"context"
	"io"

	"github.com/HumXC/shiroko/protos/common"
	pShell "github.com/HumXC/shiroko/protos/shell"
	"google.golang.org/grpc"
)

type shellClient struct {
	sc  pShell.ShellClient
	ctx context.Context
}

// Getprop implements shell.IShell.
func (s *shellClient) Getprop(key string) (string, error) {
	resp, err := s.sc.Getprop(s.ctx, &pShell.GetpropRequest{Key: key})
	if err != nil {
		return "", ParseError(err)
	}
	return resp.Value, nil
}

// Install implements shell.IShell.
func (s *shellClient) Install(apkpath string) error {
	_, err := s.sc.Install(s.ctx, &pShell.InstallRequest{Apkpath: apkpath})
	return ParseError(err)
}

// ListApps implements shell.IShell.
func (s *shellClient) ListApps() ([]string, error) {
	resp, err := s.sc.ListApps(s.ctx, &common.Empty{})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Apps, nil
}

// Pull implements shell.IShell.
func (s *shellClient) Pull(filename string) (io.ReadCloser, error) {
	resp, err := s.sc.Pull(s.ctx, &pShell.PullRequest{Filename: filename})
	if err != nil {
		return nil, ParseError(err)
	}
	return common.NewReadCloser(resp), nil
}

// Push implements shell.IShell.
func (s *shellClient) Push(filename string, data io.Reader) error {
	isSendFilename := false
	c, err := s.sc.Push(s.ctx)
	if err != nil {
		return ParseError(err)
	}
	defer c.CloseAndRecv()

	buf := make([]byte, 256)
	for {
		n, err := data.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return ParseError(err)
		}
		req := &pShell.PushRequest{Data: buf[:n]}
		if !isSendFilename {
			req.Filename = filename
			isSendFilename = true
		}
		err = c.Send(req)
		if err != nil {
			return ParseError(err)
		}
	}
	return c.CloseSend()
}

// Run implements shell.IShell.
func (s *shellClient) Run(cmd string, timeoutMs int32) ([]byte, error) {
	resp, err := s.sc.Run(s.ctx, &pShell.RunRequest{
		Cmd:       cmd,
		TimeoutMs: timeoutMs,
	})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Data, nil
}

// StartApp implements shell.IShell.
func (s *shellClient) StartApp(active string) error {
	_, err := s.sc.StartApp(s.ctx, &pShell.StartAppRequest{Active: active})
	return ParseError(err)
}

// StopApp implements shell.IShell.
func (s *shellClient) StopApp(pkgname string) error {
	_, err := s.sc.StopApp(s.ctx, &pShell.StopAppRequest{Pkgname: pkgname})
	return ParseError(err)
}

// Uninstall implements shell.IShell.
func (s *shellClient) Uninstall(pkgname string) error {
	_, err := s.sc.Uninstall(s.ctx, &pShell.UninstallRequest{Pkgname: pkgname})
	return ParseError(err)
}

func initShell(ctx context.Context, conn *grpc.ClientConn) Shell {
	return &shellClient{
		sc:  pShell.NewShellClient(conn),
		ctx: ctx,
	}
}
