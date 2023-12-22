package server

import (
	"context"

	"github.com/HumXC/shiroko/protos/common"
	pWindow "github.com/HumXC/shiroko/protos/window"
	"github.com/HumXC/shiroko/tools/window"
)

type serverWindow struct {
	// 嵌入此是为了确保所有定义的方法都被实现
	pWindow.UnsafeWindowServer
	window window.IWindow
}

// GetDensity implements window.WindowServer.
func (s *serverWindow) GetDensity(ctx context.Context, req *common.Empty) (*pWindow.Density, error) {
	d, err := s.window.GetDensity()
	if err != nil {
		return nil, err
	}
	return &pWindow.Density{Density: d}, nil
}

// GetSize implements window.WindowServer.
func (s *serverWindow) GetSize(ctx context.Context, req *common.Empty) (*pWindow.Size, error) {
	w, h, err := s.window.GetSize()
	if err != nil {
		return nil, err
	}
	return &pWindow.Size{Width: w, Height: h}, nil
}

// ResetDensity implements window.WindowServer.
func (s *serverWindow) ResetDensity(ctx context.Context, req *common.Empty) (*common.Empty, error) {
	err := s.window.ResetDensity()
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

// ResetSize implements window.WindowServer.
func (s *serverWindow) ResetSize(ctx context.Context, req *common.Empty) (*common.Empty, error) {
	err := s.window.ResetSize()
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

// SetDensity implements window.WindowServer.
func (s *serverWindow) SetDensity(ctx context.Context, req *pWindow.Density) (*common.Empty, error) {
	err := s.window.SetDensity(req.Density)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

// SetRotation implements window.WindowServer.
func (s *serverWindow) SetRotation(ctx context.Context, req *pWindow.Rotation) (*common.Empty, error) {
	err := s.window.SetRotation(req.Lock, req.Rotation)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

// SetSize implements window.WindowServer.
func (s *serverWindow) SetSize(ctx context.Context, req *pWindow.Size) (*common.Empty, error) {
	err := s.window.SetSize(req.Width, req.Height)
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

var _ pWindow.WindowServer = &serverWindow{}

func NewWindowServer() *serverWindow {
	return &serverWindow{
		window: window.Window,
	}
}
