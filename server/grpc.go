package server

import (
	pManager "github.com/HumXC/shiroko/protos/manager"
	pMinicap "github.com/HumXC/shiroko/protos/minicap"
	pScreencap "github.com/HumXC/shiroko/protos/screencap"
	"google.golang.org/grpc"
)

type Server struct {
	*serverScreencap
}

func New() *grpc.Server {
	grpcServer := grpc.NewServer()
	pManager.RegisterManagerServer(grpcServer, NewManagerServer())
	pScreencap.RegisterScreencapServer(grpcServer, NewScreencapServer())
	pMinicap.RegisterMinicapServer(grpcServer, NewMinicapServer())

	return grpcServer
}
