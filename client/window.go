package client

import (
	"context"

	"github.com/HumXC/shiroko/protos/common"
	pWindow "github.com/HumXC/shiroko/protos/window"
	"google.golang.org/grpc"
)

type windowClient struct {
	wc  pWindow.WindowClient
	ctx context.Context
}

// GetDensity implements window.IWindow.
func (w *windowClient) GetDensity() (int32, error) {
	reqp, err := w.wc.GetDensity(w.ctx, &common.Empty{})
	if err != nil {
		return 0, ParseError(err)
	}
	return reqp.GetDensity(), nil
}

// GetSize implements window.IWindow.
func (w *windowClient) GetSize() (int32, int32, error) {
	reqp, err := w.wc.GetSize(w.ctx, &common.Empty{})
	if err != nil {
		return 0, 0, ParseError(err)
	}
	return reqp.GetWidth(), reqp.GetHeight(), nil
}

// ResetDensity implements window.IWindow.
func (w *windowClient) ResetDensity() error {
	_, err := w.wc.ResetDensity(w.ctx, &common.Empty{})
	return ParseError(err)
}

// ResetSize implements window.IWindow.
func (w *windowClient) ResetSize() error {
	_, err := w.wc.ResetSize(w.ctx, &common.Empty{})
	return ParseError(err)
}

// SetDensity implements window.IWindow.
func (w *windowClient) SetDensity(density int32) error {
	_, err := w.wc.SetDensity(w.ctx, &pWindow.Density{Density: density})
	return ParseError(err)
}

// SetRotation implements window.IWindow.
func (w *windowClient) SetRotation(lock bool, rotation int32) error {
	_, err := w.wc.SetRotation(w.ctx, &pWindow.Rotation{Lock: lock, Rotation: rotation})
	return ParseError(err)
}

// SetSize implements window.IWindow.
func (w *windowClient) SetSize(width int32, height int32) error {
	_, err := w.wc.SetSize(w.ctx, &pWindow.Size{Width: width, Height: height})
	return ParseError(err)
}

func initWindow(ctx context.Context, conn *grpc.ClientConn) Window {
	return &windowClient{
		wc:  pWindow.NewWindowClient(conn),
		ctx: ctx,
	}
}
