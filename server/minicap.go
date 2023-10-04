package server

import (
	"context"
	"io"

	"github.com/HumXC/shiroko/protos/common"
	pMinicap "github.com/HumXC/shiroko/protos/minicap"
	"github.com/HumXC/shiroko/tools/minicap"
)

type serverMinicap struct {
	// 嵌入此是为了确保所有定义的方法都被实现
	pMinicap.UnsafeMinicapServer
	minicap minicap.IMinicap
}

// Jpg implements minicap.MinicapServer.
func (s *serverMinicap) Jpg(ctx context.Context, req *pMinicap.JpgRequest) (*common.DataChunk, error) {
	data, err := s.minicap.Jpg(req.RWidth, req.RHeight, req.VWidth, req.VHeight, req.Orientation, req.Quality)
	if err != nil {
		return &common.DataChunk{}, err
	}
	return &common.DataChunk{Data: data}, nil
}

// Cat implements minicap.MinicapServer.
func (s *serverMinicap) Cat(e *common.Empty, cat pMinicap.Minicap_CatServer) error {
	reader, err := s.minicap.Cat()
	if err != nil {
		return err
	}
	_, _ = io.Copy(common.NewWriter(cat), reader)
	return nil
}

// Info implements minicap.MinicapServer.
func (s *serverMinicap) Info(context.Context, *common.Empty) (*pMinicap.InfoResponse, error) {
	result, err := s.minicap.Info()
	if err != nil {
		return &pMinicap.InfoResponse{}, err
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
func (s *serverMinicap) Start(ctx context.Context, req *pMinicap.StartRequest) (*common.Empty, error) {
	err := s.minicap.Start(req.RWidth, req.RHeight, req.VWidth, req.VHeight, req.Orientation, req.Rate)
	if err != nil {
		return &common.Empty{}, err
	}
	return &common.Empty{}, nil
}

// Stop implements minicap.MinicapServer.
func (s *serverMinicap) Stop(context.Context, *common.Empty) (*common.Empty, error) {
	err := s.minicap.Stop()
	if err != nil {
		return &common.Empty{}, err
	}
	return &common.Empty{}, nil
}

var _ pMinicap.MinicapServer = &serverMinicap{}

func NewMinicapServer() *serverMinicap {
	return &serverMinicap{
		minicap: minicap.Minicap,
	}
}
