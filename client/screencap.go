package client

import (
	"context"

	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	"google.golang.org/grpc"
)

type screencapClient struct {
	sc  pScreencap.ScreencapClient
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
func (s *screencapClient) Png(displayID string) ([]byte, error) {
	resp, err := s.sc.Png(s.ctx, &pScreencap.PngRequest{
		DisplayID: displayID,
	})
	if err != nil {
		return nil, ParseError(err)
	}
	return resp.Data, nil
}

func initScreencap(ctx context.Context, conn *grpc.ClientConn) Screencap {
	return &screencapClient{
		sc:  pScreencap.NewScreencapClient(conn),
		ctx: ctx,
	}
}
