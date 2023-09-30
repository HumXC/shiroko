package server

import (
	pScreenca "github.com/HumXC/shiroko/protos/screencap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MakeError(description string, err error) error {
	return status.Errorf(codes.Aborted, description+": %w", err)
}

type Server struct {
	*serverScreencap
}

func New() *grpc.Server {
	grpcServer := grpc.NewServer()
	pScreenca.RegisterScreencapServiceServer(grpcServer, NewScreencapServer())
	return grpcServer
}
