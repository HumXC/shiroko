package server

import (
	"context"

	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	"github.com/HumXC/shiroko/tools/screencap"
)

type serverScreencap struct {
	pScreencap.UnimplementedScreencapServiceServer
	screencap *screencap.Screencap
}

func NewScreencapServer() *serverScreencap {
	return &serverScreencap{
		screencap: screencap.New(),
	}
}
func (s *serverScreencap) Displays(ctx context.Context, req *pScreencap.DisplaysRequest) (*pScreencap.DisplaysResponse, error) {
	result, err := s.screencap.Displays()
	if err != nil {
		return nil, MakeError("failed to get displays", err)
	}
	return &pScreencap.DisplaysResponse{
		DisplayIDs: result,
	}, nil
}

// Png implements protos.ScreencapServiceServer.
func (s *serverScreencap) Png(context.Context, *pScreencap.PngRequest) (*pScreencap.PngResponse, error) {
	result, err := s.screencap.Png()
	if err != nil {
		return nil, MakeError("failed to screencap", err)
	}
	return &pScreencap.PngResponse{
		Data: result,
	}, nil
}

// PngWithDisplay implements protos.ScreencapServiceServer.
func (s *serverScreencap) PngWithDisplay(ctx context.Context, req *pScreencap.PngWithDisplayRequest) (*pScreencap.PngResponse, error) {
	result, err := s.screencap.PngWithDisplay(req.DisplayID)
	if err != nil {
		return nil, MakeError("failed to screencap with id", err)
	}
	return &pScreencap.PngResponse{
		Data: result,
	}, nil
}
