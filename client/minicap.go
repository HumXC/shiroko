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
type grpcCatReadCloser struct {
	stream pMinicap.Minicap_CatClient
	buf    []byte
}

// Close implements io.ReadCloser.
func (g *grpcCatReadCloser) Close() error {
	return g.stream.CloseSend()
}

func (g *grpcCatReadCloser) Read(p []byte) (n int, err error) {
	if len(g.buf) == 0 { // 如果buffer中没有数据，尝试从流中获取
		dataChunk, err := g.stream.Recv()
		if err == io.EOF {
			return 0, io.EOF
		}
		if err != nil {
			return 0, err
		}
		g.buf = dataChunk.Data
	}
	// 从buffer中复制数据到p
	n = copy(p, g.buf)
	g.buf = g.buf[n:]
	return n, nil
}

// Cat implements minicap.IMinicap.
func (m *minicapClient) Cat() (io.ReadCloser, error) {
	catClient, err := m.mm.Cat(m.ctx, &pMinicap.Empty{})
	if err != nil {
		return nil, ParseError(err)
	}
	return &grpcCatReadCloser{
		stream: catClient,
	}, nil
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
func (m *minicapClient) Start(RWidth, RHeight, VWidth, VHeight, Orientation int32) error {
	_, err := m.mm.Start(m.ctx, &pMinicap.StartRequest{
		RWidth:      RWidth,
		RHeight:     RHeight,
		VWidth:      VWidth,
		VHeight:     VHeight,
		Orientation: Orientation})
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
