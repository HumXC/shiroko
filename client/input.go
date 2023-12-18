package client

import (
	"context"

	pInput "github.com/HumXC/shiroko/protos/input"
	"google.golang.org/grpc"
)

type inputClient struct {
	ic  pInput.InputClient
	ctx context.Context
}

// Draganddrop implements input.IInput.
func (i *inputClient) Draganddrop(x1 int32, y1 int32, x2 int32, y2 int32, duration int32) error {
	_, err := i.ic.Draganddrop(i.ctx, &pInput.DraganddropRequest{
		X1:       x1,
		Y1:       y1,
		X2:       x2,
		Y2:       y2,
		Duration: duration,
	})
	return ParseError(err)
}

// Keyevent implements input.IInput.
func (i *inputClient) Keyevent(longpress bool, keycode ...string) error {
	_, err := i.ic.Keyevent(i.ctx, &pInput.KeyeventRequest{
		Longpress: longpress,
		Keycode:   keycode,
	})
	return ParseError(err)
}

// Motionevent implements input.IInput.
func (i *inputClient) Motionevent(event string, x int32, y int32) error {
	_, err := i.ic.Motionevent(i.ctx, &pInput.MotioneventRequest{
		Event: event,
		X:     x,
		Y:     y,
	})
	return ParseError(err)
}

// Press implements input.IInput.
func (i *inputClient) Press() error {
	_, err := i.ic.Press(i.ctx, &pInput.PressRequest{})
	return ParseError(err)
}

// Roll implements input.IInput.
func (i *inputClient) Roll(dx int32, dy int32) error {
	_, err := i.ic.Roll(i.ctx, &pInput.RollRequest{
		Dx: dx,
		Dy: dy,
	})
	return ParseError(err)
}

// Swipe implements input.IInput.
func (i *inputClient) Swipe(x1 int32, y1 int32, x2 int32, y2 int32, duration int32) error {
	_, err := i.ic.Swipe(i.ctx, &pInput.SwipeRequest{
		X1:       x1,
		Y1:       y1,
		X2:       x2,
		Y2:       y2,
		Duration: duration,
	})
	return ParseError(err)
}

// Tap implements input.IInput.
func (i *inputClient) Tap(x int32, y int32) error {
	_, err := i.ic.Tap(i.ctx, &pInput.TapRequest{
		X: x,
		Y: y,
	})
	return ParseError(err)
}

// Text implements input.IInput.
func (i *inputClient) Text(text string) error {
	_, err := i.ic.Text(i.ctx, &pInput.TextRequest{
		Text: text,
	})
	return ParseError(err)
}

func initInput(ctx context.Context, conn *grpc.ClientConn) Input {
	return &inputClient{
		ic:  pInput.NewInputClient(conn),
		ctx: ctx,
	}
}
