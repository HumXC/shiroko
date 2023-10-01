package server

import (
	"context"
	"io"

	pMinicap "github.com/HumXC/shiroko/protos/minicap"
	"github.com/HumXC/shiroko/tools/minicap"
)

type serverMinicap struct {
	pMinicap.UnimplementedMinicapServer
	minicap minicap.IMinicap
}

// Jpg implements minicap.MinicapServer.
func (s *serverMinicap) Jpg(ctx context.Context, req *pMinicap.JpgRequest) (*pMinicap.JpgResponse, error) {
	data, err := s.minicap.Jpg(req.RWidth, req.RHeight, req.VWidth, req.VHeight, req.Orientation, req.Quality)
	if err != nil {
		return &pMinicap.JpgResponse{}, MakeError("failed to start minicap", err)
	}
	return &pMinicap.JpgResponse{Data: data}, nil
}

type grpcMinicapWriter struct {
	stream pMinicap.Minicap_CatServer
}

func (w *grpcMinicapWriter) Write(p []byte) (n int, err error) {
	chunk := &pMinicap.DataChunk{Data: p}
	if err := w.stream.Send(chunk); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Cat implements minicap.MinicapServer.
func (s *serverMinicap) Cat(e *pMinicap.Empty, cat pMinicap.Minicap_CatServer) error {
	reader, err := s.minicap.Cat()
	if err != nil {
		return MakeError("failed to get minicap socket", err)
	}
	_, _ = io.Copy(&grpcMinicapWriter{cat}, reader)
	return nil
}

// Info implements minicap.MinicapServer.
func (s *serverMinicap) Info(context.Context, *pMinicap.Empty) (*pMinicap.InfoResponse, error) {
	result, err := s.minicap.Info()
	if err != nil {
		return &pMinicap.InfoResponse{}, MakeError("failed to get minicap info", err)
	}
	return &pMinicap.InfoResponse{
		Id:       result.Id,
		Width:    result.Width,
		Height:   result.Height,
		Xdpi:     result.Xdpi,
		Ydpi:     result.Ydpi,
		Size:     result.Size,
		Density:  result.Density,
		Fps:      result.Fps,
		Secure:   result.Secure,
		Rotation: result.Rotation,
	}, nil
}

// Start implements minicap.MinicapServer.
func (s *serverMinicap) Start(ctx context.Context, req *pMinicap.StartRequest) (*pMinicap.Empty, error) {
	err := s.minicap.Start(req.RWidth, req.RHeight, req.VWidth, req.VHeight, req.Orientation, req.Rate)
	if err != nil {
		return &pMinicap.Empty{}, MakeError("failed to start minicap", err)
	}
	return &pMinicap.Empty{}, nil
}

// Stop implements minicap.MinicapServer.
func (s *serverMinicap) Stop(context.Context, *pMinicap.Empty) (*pMinicap.Empty, error) {
	err := s.minicap.Stop()
	if err != nil {
		return &pMinicap.Empty{}, MakeError("failed to stop minicap", err)
	}
	return &pMinicap.Empty{}, nil
}

var _ pMinicap.MinicapServer = &serverMinicap{}

func NewMinicapServer() *serverMinicap {
	return &serverMinicap{
		minicap: minicap.Minicap,
	}
}
