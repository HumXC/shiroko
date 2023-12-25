package server

import (
	"context"

	"github.com/HumXC/shiroko/protos/common"
	pInput "github.com/HumXC/shiroko/protos/input"
	"github.com/HumXC/shiroko/tools/input"
)

type serverInput struct {
	// 嵌入此是为了确保所有定义的方法都被实现
	pInput.UnsafeInputServer
	input input.IInput
}

// Draganddrop implements input.InputServer.
func (s *serverInput) Draganddrop(ctx context.Context, req *pInput.DraganddropRequest) (*common.Empty, error) {
	err := s.input.Draganddrop(req.X1, req.Y1, req.X2, req.Y2, req.Duration)
	return &common.Empty{}, err
}

// Keyevent implements input.InputServer.
func (s *serverInput) Keyevent(ctx context.Context, req *pInput.KeyeventRequest) (*common.Empty, error) {
	err := s.input.Keyevent(req.Longpress, req.Keycode...)
	return &common.Empty{}, err
}

// Motionevent implements input.InputServer.
func (s *serverInput) Motionevent(ctx context.Context, req *pInput.MotioneventRequest) (*common.Empty, error) {
	err := s.input.Motionevent(req.Event, req.X, req.Y)
	return &common.Empty{}, err
}

// Swipe implements input.InputServer.
func (s *serverInput) Swipe(ctx context.Context, req *pInput.SwipeRequest) (*common.Empty, error) {
	err := s.input.Swipe(req.X1, req.Y1, req.X2, req.Y2, req.Duration)
	return &common.Empty{}, err
}

// Tap implements input.InputServer.
func (s *serverInput) Tap(ctx context.Context, req *pInput.TapRequest) (*common.Empty, error) {
	err := s.input.Tap(req.X, req.Y)
	return &common.Empty{}, err
}

// Text implements input.InputServer.
func (s *serverInput) Text(ctx context.Context, req *pInput.TextRequest) (*common.Empty, error) {
	err := s.input.Text(req.Text)
	return &common.Empty{}, err
}

var _ pInput.InputServer = &serverInput{}

func NewInputServer() *serverInput {
	return &serverInput{
		input: input.Input,
	}
}
