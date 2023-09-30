package client

import (
	"context"

	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	tScreencap "github.com/HumXC/shiroko/tools/screencap"
	"google.golang.org/grpc"
)

type screencapClient struct {
	sc  pScreencap.ScreencapServiceClient
	ctx context.Context
}

// Displays implements screencap.IScreencap.
func (s *screencapClient) Displays() ([]string, error) {
	resp, err := s.sc.Displays(s.ctx, &pScreencap.DisplaysRequest{})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.DisplayIDs, nil
}

// Png implements screencap.IScreencap.
func (s *screencapClient) Png() ([]byte, error) {
	resp, err := s.sc.Png(s.ctx, &pScreencap.PngRequest{})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Data, nil
}

// PngWithDisplay implements screencap.IScreencap.
func (s *screencapClient) PngWithDisplay(displayID string) ([]byte, error) {
	resp, err := s.sc.PngWithDisplay(s.ctx, &pScreencap.PngWithDisplayRequest{DisplayID: displayID})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Data, nil
}

func initScreencap(ctx context.Context, conn *grpc.ClientConn) tScreencap.IScreencap {
	return &screencapClient{
		sc:  pScreencap.NewScreencapServiceClient(conn),
		ctx: ctx,
	}
}
