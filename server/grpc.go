package server

import (
	"io"

	pMinicap "github.com/HumXC/shiroko/protos/minicap"
	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MakeError(description string, err error) error {
	return status.Errorf(codes.Aborted, description+": %s", err)
}

type Server struct {
	*serverScreencap
}

func New() *grpc.Server {
	grpcServer := grpc.NewServer()
	pScreencap.RegisterScreencapServer(grpcServer, NewScreencapServer())
	pMinicap.RegisterMinicapServer(grpcServer, NewMinicapServer())
	return grpcServer
}

type writer struct {
	stream grpc.ServerStream
}

func (w *writer) Write(p []byte) (n int, err error) {
	chunk := &pMinicap.DataChunk{Data: p}
	if err := w.stream.SendMsg(chunk); err != nil {
		return 0, err
	}
	return len(p), nil
}
func NewWriter(stream grpc.ServerStream) io.Writer {
	return &writer{stream}
}
