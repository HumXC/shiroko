package server

import (
	"context"

	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	"github.com/HumXC/shiroko/tools/screencap"
)

type serverScreencap struct {
	pScreencap.UnimplementedScreencapServiceServer
	screencap screencap.IScreencap
}

var _ pScreencap.ScreencapServiceServer = &serverScreencap{}

// Png implements screencap.ScreencapServiceServer.
func (s *serverScreencap) Png(ctx context.Context, req *pScreencap.PngRequest) (*pScreencap.PngResponse, error) {
	result, err := s.screencap.Png(req.DisplayID)
	if err != nil {
		return nil, MakeError("failed to screencap", err)
	}
	return &pScreencap.PngResponse{
		Data: result,
	}, nil
}

// Displays implements screencap.ScreencapServiceServer.
func (s *serverScreencap) Displays(ctx context.Context, req *pScreencap.DisplaysRequest) (*pScreencap.DisplaysResponse, error) {
	result, err := s.screencap.Displays()
	if err != nil {
		return nil, MakeError("failed to get displays", err)
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
