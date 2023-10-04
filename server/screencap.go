package server

import (
	"context"

	"github.com/HumXC/shiroko/protos/common"
	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	"github.com/HumXC/shiroko/tools/screencap"
)

type serverScreencap struct {
	// 嵌入此是为了确保所有定义的方法都被实现
	pScreencap.UnsafeScreencapServer
	screencap screencap.IScreencap
}

// Png implements screencap.ScreencapServer.
func (s *serverScreencap) Png(ctx context.Context, req *pScreencap.PngRequest) (*common.DataChunk, error) {
	result, err := s.screencap.Png(req.DisplayID)
	if err != nil {
		return nil, err
	}
	return &common.DataChunk{Data: result}, nil
}

var _ pScreencap.ScreencapServer = &serverScreencap{}

// Displays implements screencap.ScreencapServiceServer.
func (s *serverScreencap) Displays(ctx context.Context, req *pScreencap.DisplaysRequest) (*pScreencap.DisplaysResponse, error) {
	result, err := s.screencap.Displays()
	if err != nil {
		return nil, err
	}
	return &pScreencap.DisplaysResponse{
		DisplayIDs: result,
	}, nil
}

func NewScreencapServer() *serverScreencap {
	return &serverScreencap{
		screencap: screencap.Screencap,
	}
}
