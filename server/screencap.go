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
func (s *serverScreencap) Png(ctx context.Context, req *common.Empty) (*common.DataChunk, error) {
	result, err := s.screencap.Png()
	if err != nil {
		return nil, err
	}
	return &common.DataChunk{Data: result}, nil
}

var _ pScreencap.ScreencapServer = &serverScreencap{}

func NewScreencapServer() *serverScreencap {
	return &serverScreencap{
		screencap: screencap.Screencap,
	}
}
