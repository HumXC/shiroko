package client

import (
	"context"
	"io"

	pMinicap "github.com/HumXC/shiroko/protos/minicap"
	tMinicap "github.com/HumXC/shiroko/tools/minicap"
	"google.golang.org/grpc"
)

type minicapClient struct {
	mm  pMinicap.MinicapClient
	ctx context.Context
}

// Jpg implements minicap.IMinicap.
func (m *minicapClient) Jpg(rWidth int32, rHeight int32, vWidth int32, vHeight int32, orientation int32, quality int32) ([]byte, error) {
	resp, err := m.mm.Jpg(m.ctx, &pMinicap.JpgRequest{
		RWidth:      rWidth,
		RHeight:     rHeight,
		VWidth:      vWidth,
		VHeight:     vHeight,
		Orientation: orientation,
		Quality:     quality,
	})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Data, nil
}

// Cat implements minicap.IMinicap.
func (m *minicapClient) Cat() (io.ReadCloser, error) {
	catClient, err := m.mm.Cat(m.ctx, &pMinicap.Empty{})
	if err != nil {
		return nil, ParseError(err)
	}
	return NewReadCloser(catClient, &pMinicap.DataChunk{}), nil
}

// Info implements minicap.IMinicap.
func (m *minicapClient) Info() (tMinicap.Info, error) {
	resp, err := m.mm.Info(m.ctx, &pMinicap.Empty{})
	if err != nil {
		return tMinicap.Info{}, ParseError(err)
	}
	return tMinicap.Info{
		Id:       resp.Id,
		Width:    resp.Width,
		Height:   resp.Height,
		Xdpi:     resp.Xdpi,
		Ydpi:     resp.Ydpi,
		Size:     resp.Size,
		Density:  resp.Density,
		Fps:      resp.Fps,
		Secure:   resp.Secure,
		Rotation: resp.Rotation,
	}, nil
}

// Start implements minicap.IMinicap.
func (m *minicapClient) Start(rWidth, rHeight, vWidth, vHeight, orientation, rate int32) error {
	_, err := m.mm.Start(m.ctx, &pMinicap.StartRequest{
		RWidth:      rWidth,
		RHeight:     rHeight,
		VWidth:      vWidth,
		VHeight:     vHeight,
		Orientation: orientation,
		Rate:        rate,
	})
	if err != nil {
		return ParseError(err)
	}
	return nil
}

// Stop implements minicap.IMinicap.
func (m *minicapClient) Stop() error {
	_, err := m.mm.Stop(m.ctx, &pMinicap.Empty{})
	if err != nil {
		return ParseError(err)
	}
	return nil
}

func initMinicap(ctx context.Context, conn *grpc.ClientConn) tMinicap.IMinicap {
	return &minicapClient{
		mm:  pMinicap.NewMinicapClient(conn),
		ctx: ctx,
	}
}
