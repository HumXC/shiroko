package client

import (
	"context"

	"github.com/HumXC/shiroko/protos/common"
	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	"google.golang.org/grpc"
)

type screencapClient struct {
	sc  pScreencap.ScreencapClient
	ctx context.Context
}

// Png implements screencap.IScreencap.
func (s *screencapClient) Png() ([]byte, error) {
	resp, err := s.sc.Png(s.ctx, &common.Empty{})
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
